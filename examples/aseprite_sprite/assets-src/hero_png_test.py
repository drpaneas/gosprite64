from pathlib import Path
import pytest
import shutil
import struct
import subprocess
import zlib


def png_size(path: Path) -> tuple[int, int]:
    with path.open("rb") as png_file:
        assert png_file.read(8) == b"\x89PNG\r\n\x1a\n"
        ihdr_length = struct.unpack(">I", png_file.read(4))[0]
        assert ihdr_length == 13
        assert png_file.read(4) == b"IHDR"
        width, height = struct.unpack(">II", png_file.read(8))
        return width, height


def png_has_transparent_pixels(path: Path) -> bool:
    _, _, rgba = png_rgba_data(path)
    return any(rgba[i] < 255 for i in range(3, len(rgba), 4))


def _paeth(a: int, b: int, c: int) -> int:
    p = a + b - c
    pa = abs(p - a)
    pb = abs(p - b)
    pc = abs(p - c)
    if pa <= pb and pa <= pc:
        return a
    if pb <= pc:
        return b
    return c


def aseprite_info(path: Path) -> tuple[int, int, int]:
    with path.open("rb") as aseprite_file:
        header = aseprite_file.read(12)

    assert len(header) == 12, "expected Aseprite header"

    magic = struct.unpack("<H", header[4:6])[0]
    frames = struct.unpack("<H", header[6:8])[0]
    width = struct.unpack("<H", header[8:10])[0]
    height = struct.unpack("<H", header[10:12])[0]

    assert magic == 0xA5E0, "expected Aseprite file magic"
    return width, height, frames


def exported_png_matches(source: Path, committed_png: Path, tmp_dir: Path) -> bool:
    export_path = tmp_dir / "hero-export.png"
    aseprite_cli = find_aseprite_cli()
    if not aseprite_cli:
        pytest.skip("Aseprite CLI is unavailable; skipping source-to-export parity check")

    result = subprocess.run(
        [aseprite_cli, "-b", str(source), "--save-as", str(export_path)],
        capture_output=True,
        text=True,
        check=False,
    )
    assert result.returncode == 0, result.stderr or result.stdout
    assert export_path.exists(), f"expected exported sprite at {export_path}"

    return png_rgba_data(committed_png) == png_rgba_data(export_path)


def find_aseprite_cli() -> str | None:
    preferred = Path("/Applications/Aseprite.app/Contents/MacOS/aseprite")
    if preferred.exists():
        return str(preferred)

    return shutil.which("aseprite")


def png_rgba_data(path: Path) -> tuple[int, int, bytes]:
    with path.open("rb") as png_file:
        assert png_file.read(8) == b"\x89PNG\r\n\x1a\n"

        width = height = color_type = None
        idat_data = bytearray()

        while True:
            chunk_length = struct.unpack(">I", png_file.read(4))[0]
            chunk_type = png_file.read(4)
            chunk_data = png_file.read(chunk_length)
            png_file.read(4)

            if chunk_type == b"IHDR":
                width, height = struct.unpack(">II", chunk_data[:8])
                color_type = chunk_data[9]
                assert chunk_data[12] == 0, "interlaced PNGs are not expected here"
            elif chunk_type == b"IDAT":
                idat_data.extend(chunk_data)
            elif chunk_type == b"IEND":
                break

    assert width is not None and height is not None
    assert color_type == 6, "expected RGBA export for transparent sprite"

    raw = zlib.decompress(bytes(idat_data))
    stride = width * 4
    prev_row = bytes(stride)
    offset = 0
    rgba_rows = bytearray()

    for _ in range(height):
        filter_type = raw[offset]
        offset += 1
        filtered = bytearray(raw[offset : offset + stride])
        offset += stride

        if filter_type == 1:
            for i in range(4, stride):
                filtered[i] = (filtered[i] + filtered[i - 4]) & 0xFF
        elif filter_type == 2:
            for i in range(stride):
                filtered[i] = (filtered[i] + prev_row[i]) & 0xFF
        elif filter_type == 3:
            for i in range(stride):
                left = filtered[i - 4] if i >= 4 else 0
                up = prev_row[i]
                filtered[i] = (filtered[i] + ((left + up) // 2)) & 0xFF
        elif filter_type == 4:
            for i in range(stride):
                left = filtered[i - 4] if i >= 4 else 0
                up = prev_row[i]
                up_left = prev_row[i - 4] if i >= 4 else 0
                filtered[i] = (filtered[i] + _paeth(left, up, up_left)) & 0xFF
        else:
            assert filter_type == 0, f"unexpected PNG filter type {filter_type}"

        rgba_rows.extend(filtered)
        prev_row = bytes(filtered)

    return width, height, bytes(rgba_rows)


def test_hero_png_exists_and_is_64x64() -> None:
    hero_png = Path(__file__).with_name("hero.png")

    assert hero_png.exists(), f"expected exported sprite at {hero_png}"
    assert png_size(hero_png) == (64, 64)


def test_hero_png_contains_transparency() -> None:
    hero_png = Path(__file__).with_name("hero.png")

    assert hero_png.exists(), f"expected exported sprite at {hero_png}"
    assert png_has_transparent_pixels(hero_png)


def test_hero_aseprite_source_is_single_frame_64x64() -> None:
    hero_aseprite = Path(__file__).with_name("hero.aseprite")

    assert hero_aseprite.exists(), f"expected Aseprite source at {hero_aseprite}"
    assert aseprite_info(hero_aseprite) == (64, 64, 1)


def test_hero_png_matches_fresh_aseprite_export(tmp_path: Path) -> None:
    hero_aseprite = Path(__file__).with_name("hero.aseprite")
    hero_png = Path(__file__).with_name("hero.png")

    assert hero_aseprite.exists(), f"expected Aseprite source at {hero_aseprite}"
    assert hero_png.exists(), f"expected exported sprite at {hero_png}"
    assert exported_png_matches(hero_aseprite, hero_png, tmp_path)


def test_fresh_export_check_skips_without_aseprite_cli(
    monkeypatch: pytest.MonkeyPatch, tmp_path: Path
) -> None:
    hero_aseprite = Path(__file__).with_name("hero.aseprite")
    hero_png = Path(__file__).with_name("hero.png")
    preferred = Path("/Applications/Aseprite.app/Contents/MacOS/aseprite")
    original_exists = Path.exists

    def fake_exists(self: Path) -> bool:
        if self == preferred:
            return False
        return original_exists(self)

    monkeypatch.setattr(Path, "exists", fake_exists)
    monkeypatch.setattr(shutil, "which", lambda _: None)

    with pytest.raises(pytest.skip.Exception, match="Aseprite CLI is unavailable"):
        exported_png_matches(hero_aseprite, hero_png, tmp_path)
