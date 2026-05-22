from __future__ import annotations

import struct
import zlib
from pathlib import Path


ASSET_PATH = Path(__file__).with_name("slime_right_32.png")
EXPECTED_SIZE = (32, 32)
PNG_SIGNATURE = b"\x89PNG\r\n\x1a\n"


def _read_chunks(data: bytes) -> list[tuple[bytes, bytes]]:
    if not data.startswith(PNG_SIGNATURE):
        raise AssertionError(f"{ASSET_PATH.name} is not a PNG file")

    chunks: list[tuple[bytes, bytes]] = []
    offset = len(PNG_SIGNATURE)
    while offset < len(data):
        if offset + 8 > len(data):
            raise AssertionError("truncated PNG chunk header")
        length = struct.unpack(">I", data[offset : offset + 4])[0]
        chunk_type = data[offset + 4 : offset + 8]
        chunk_start = offset + 8
        chunk_end = chunk_start + length
        crc_end = chunk_end + 4
        if crc_end > len(data):
            raise AssertionError(f"truncated PNG chunk {chunk_type!r}")
        chunks.append((chunk_type, data[chunk_start:chunk_end]))
        offset = crc_end
        if chunk_type == b"IEND":
            break
    return chunks


def _bytes_per_pixel(color_type: int) -> int:
    return {
        0: 1,
        2: 3,
        3: 1,
        4: 2,
        6: 4,
    }[color_type]


def _unfilter_scanlines(raw: bytes, width: int, height: int, color_type: int) -> list[bytes]:
    bpp = _bytes_per_pixel(color_type)
    stride = width * bpp
    expected = height * (stride + 1)
    if len(raw) != expected:
        raise AssertionError(
            f"unexpected decompressed image size: got {len(raw)}, want {expected}"
        )

    rows: list[bytes] = []
    prev = bytearray(stride)
    pos = 0

    for _ in range(height):
        filter_type = raw[pos]
        pos += 1
        row = bytearray(raw[pos : pos + stride])
        pos += stride

        if filter_type == 0:
            pass
        elif filter_type == 1:
            for i in range(stride):
                left = row[i - bpp] if i >= bpp else 0
                row[i] = (row[i] + left) & 0xFF
        elif filter_type == 2:
            for i in range(stride):
                row[i] = (row[i] + prev[i]) & 0xFF
        elif filter_type == 3:
            for i in range(stride):
                left = row[i - bpp] if i >= bpp else 0
                up = prev[i]
                row[i] = (row[i] + ((left + up) // 2)) & 0xFF
        elif filter_type == 4:
            for i in range(stride):
                left = row[i - bpp] if i >= bpp else 0
                up = prev[i]
                up_left = prev[i - bpp] if i >= bpp else 0
                predictor = _paeth(left, up, up_left)
                row[i] = (row[i] + predictor) & 0xFF
        else:
            raise AssertionError(f"unsupported PNG filter type: {filter_type}")

        rows.append(bytes(row))
        prev = row

    return rows


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


def _has_transparent_pixel(chunks: list[tuple[bytes, bytes]], width: int, height: int) -> bool:
    chunk_map = {chunk_type: payload for chunk_type, payload in chunks}
    ihdr = chunk_map[b"IHDR"]
    bit_depth = ihdr[8]
    color_type = ihdr[9]

    if bit_depth != 8:
        raise AssertionError(f"unsupported PNG bit depth: {bit_depth}")

    if color_type == 3 and b"tRNS" in chunk_map:
        return True
    if color_type not in {4, 6}:
        return b"tRNS" in chunk_map

    compressed = b"".join(payload for chunk_type, payload in chunks if chunk_type == b"IDAT")
    raw = zlib.decompress(compressed)
    rows = _unfilter_scanlines(raw, width, height, color_type)
    alpha_offset = 1 if color_type == 4 else 3
    pixel_stride = _bytes_per_pixel(color_type)

    for row in rows:
        for i in range(alpha_offset, len(row), pixel_stride):
            if row[i] < 255:
                return True
    return False


def test_selected_pixellab_source_asset_exists_and_has_expected_shape() -> None:
    assert ASSET_PATH.exists(), f"missing selected source asset: {ASSET_PATH.name}"

    data = ASSET_PATH.read_bytes()
    chunks = _read_chunks(data)
    chunk_map = {chunk_type: payload for chunk_type, payload in chunks}

    assert b"IHDR" in chunk_map, "PNG missing IHDR chunk"

    width, height = struct.unpack(">II", chunk_map[b"IHDR"][:8])
    assert (width, height) == EXPECTED_SIZE
    assert _has_transparent_pixel(chunks, width, height), (
        "expected transparent padding/background for sprite compositing"
    )
