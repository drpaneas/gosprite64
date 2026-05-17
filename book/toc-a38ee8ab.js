// Populate the sidebar
//
// This is a script, and not included directly in the page, to control the total size of the book.
// The TOC contains an entry for each page, so if each page includes a copy of the TOC,
// the total size of the page becomes O(n**2).
class MDBookSidebarScrollbox extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        this.innerHTML = '<ol class="chapter"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">1.</strong> Part 1: Welcome</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="01-welcome/why-gosprite64.html"><strong aria-hidden="true">1.1.</strong> Why GoSprite64</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="01-welcome/feature-overview.html"><strong aria-hidden="true">1.2.</strong> Feature Overview</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">2.</strong> Part 2: Getting Started</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="02-getting-started/installation.html"><strong aria-hidden="true">2.1.</strong> Installation</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="02-getting-started/hello-world.html"><strong aria-hidden="true">2.2.</strong> Hello World</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="02-getting-started/editor-setup.html"><strong aria-hidden="true">2.3.</strong> Editor Setup</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">3.</strong> Part 3: Tutorial - Build a Platformer</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/01-start-the-engine.html"><strong aria-hidden="true">3.1.</strong> Step 1: Start the Engine</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/02-draw-a-tilemap.html"><strong aria-hidden="true">3.2.</strong> Step 2: Draw a Tilemap World</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/03-add-a-player-sprite.html"><strong aria-hidden="true">3.3.</strong> Step 3: Add a Player Sprite</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/04-animate-the-player.html"><strong aria-hidden="true">3.4.</strong> Step 4: Animate the Player</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/05-move-with-dpad.html"><strong aria-hidden="true">3.5.</strong> Step 5: Move with the D-Pad</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/06-camera-following.html"><strong aria-hidden="true">3.6.</strong> Step 6: Camera Following</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/07-add-sound-effects.html"><strong aria-hidden="true">3.7.</strong> Step 7: Add Sound Effects</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/08-add-title-screen.html"><strong aria-hidden="true">3.8.</strong> Step 8: Add a Title Screen</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/09-screen-transitions.html"><strong aria-hidden="true">3.9.</strong> Step 9: Screen Transitions</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/10-score-display.html"><strong aria-hidden="true">3.10.</strong> Step 10: Score Display</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="03-tutorial/11-final-polish.html"><strong aria-hidden="true">3.11.</strong> Step 11: Final Polish</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">4.</strong> Part 4: Core Concepts</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="04-core-concepts/game-loop.html"><strong aria-hidden="true">4.1.</strong> The Game Loop</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="04-core-concepts/fixed-canvas.html"><strong aria-hidden="true">4.2.</strong> The Fixed Canvas</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="04-core-concepts/square-pixels.html"><strong aria-hidden="true">4.3.</strong> Square Pixels</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="04-core-concepts/colors.html"><strong aria-hidden="true">4.4.</strong> Colors</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">5.</strong> Part 5: Graphics</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/drawing-primitives.html"><strong aria-hidden="true">5.1.</strong> Drawing Primitives</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/sprites.html"><strong aria-hidden="true">5.2.</strong> Sprites</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/sprite-sheets.html"><strong aria-hidden="true">5.3.</strong> Sprite Sheets</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/animation-player.html"><strong aria-hidden="true">5.4.</strong> Animation Player</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/custom-fonts.html"><strong aria-hidden="true">5.5.</strong> Custom Fonts</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/text-alignment.html"><strong aria-hidden="true">5.6.</strong> Text Alignment</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/parallax.html"><strong aria-hidden="true">5.7.</strong> Parallax Scrolling</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/transitions.html"><strong aria-hidden="true">5.8.</strong> Screen Transitions</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="05-graphics/draw-regions.html"><strong aria-hidden="true">5.9.</strong> Draw Regions</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">6.</strong> Part 6: Input</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="06-input/buttons-and-dpad.html"><strong aria-hidden="true">6.1.</strong> D-Pad and Buttons</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="06-input/analog-stick.html"><strong aria-hidden="true">6.2.</strong> Analog Stick</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="06-input/multi-controller.html"><strong aria-hidden="true">6.3.</strong> Multi-Controller Support</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="06-input/rumble.html"><strong aria-hidden="true">6.4.</strong> Rumble</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="06-input/input-replay.html"><strong aria-hidden="true">6.5.</strong> Input Recording and Replay</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">7.</strong> Part 7: Audio</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="07-audio/sfx-and-music.html"><strong aria-hidden="true">7.1.</strong> Sound Effects and Music</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="07-audio/sequence-player.html"><strong aria-hidden="true">7.2.</strong> Sequence Player</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="07-audio/instrument-banks.html"><strong aria-hidden="true">7.3.</strong> Instrument Banks</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">8.</strong> Part 8: Tile Scenes</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="08-tile-scenes/pipeline-overview.html"><strong aria-hidden="true">8.1.</strong> Tile2D Pipeline Overview</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="08-tile-scenes/tile-sheets-and-maps.html"><strong aria-hidden="true">8.2.</strong> Tile Sheets and Maps</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="08-tile-scenes/bundles-and-loading.html"><strong aria-hidden="true">8.3.</strong> Bundles and Loading</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="08-tile-scenes/camera-and-scrolling.html"><strong aria-hidden="true">8.4.</strong> Camera and Scrolling</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">9.</strong> Part 9: Game Systems</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="09-game-systems/state-machine.html"><strong aria-hidden="true">9.1.</strong> State Machine</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="09-game-systems/timers.html"><strong aria-hidden="true">9.2.</strong> Timers</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="09-game-systems/menus.html"><strong aria-hidden="true">9.3.</strong> Menus</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="09-game-systems/save-data.html"><strong aria-hidden="true">9.4.</strong> Save Data</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">10.</strong> Part 10: 2D Math</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/vectors.html"><strong aria-hidden="true">10.1.</strong> Vectors</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/rectangles.html"><strong aria-hidden="true">10.2.</strong> Rectangles</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/collision-detection.html"><strong aria-hidden="true">10.3.</strong> Collision Detection</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/easing-functions.html"><strong aria-hidden="true">10.4.</strong> Easing Functions</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/grid-utilities.html"><strong aria-hidden="true">10.5.</strong> Grid Utilities</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="10-math/random-numbers.html"><strong aria-hidden="true">10.6.</strong> Random Numbers</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">11.</strong> Part 11: 3D Graphics</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="11-3d-graphics/3d-math.html"><strong aria-hidden="true">11.1.</strong> 3D Math</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="11-3d-graphics/scene-graph.html"><strong aria-hidden="true">11.2.</strong> Scene Graph</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="11-3d-graphics/display-lists.html"><strong aria-hidden="true">11.3.</strong> Display Lists</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="11-3d-graphics/triangle-rendering.html"><strong aria-hidden="true">11.4.</strong> Triangle Rendering</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">12.</strong> Part 12: Low-Level</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="12-low-level/dma-transfers.html"><strong aria-hidden="true">12.1.</strong> DMA Transfers</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="12-low-level/rsp-task-queue.html"><strong aria-hidden="true">12.2.</strong> RSP Task Queue</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="12-low-level/n64-os-primitives.html"><strong aria-hidden="true">12.3.</strong> N64 OS Primitives</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="12-low-level/memory-pools.html"><strong aria-hidden="true">12.4.</strong> Memory Pools</a></span></li></ol><li class="chapter-item expanded "><span class="chapter-link-wrapper"><span><strong aria-hidden="true">13.</strong> Part 13: Reference</span></span><ol class="section"><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="13-reference/api-quick-reference.html"><strong aria-hidden="true">13.1.</strong> API Quick Reference</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="13-reference/performance-notes.html"><strong aria-hidden="true">13.2.</strong> Performance Notes</a></span></li><li class="chapter-item expanded "><span class="chapter-link-wrapper"><a href="13-reference/troubleshooting.html"><strong aria-hidden="true">13.3.</strong> Troubleshooting</a></span></li></ol></li></ol>';
        // Set the current, active page, and reveal it if it's hidden
        let current_page = document.location.href.toString().split('#')[0].split('?')[0];
        if (current_page.endsWith('/')) {
            current_page += 'index.html';
        }
        const links = Array.prototype.slice.call(this.querySelectorAll('a'));
        const l = links.length;
        for (let i = 0; i < l; ++i) {
            const link = links[i];
            const href = link.getAttribute('href');
            if (href && !href.startsWith('#') && !/^(?:[a-z+]+:)?\/\//.test(href)) {
                link.href = path_to_root + href;
            }
            // The 'index' page is supposed to alias the first chapter in the book.
            if (link.href === current_page
                || i === 0
                && path_to_root === ''
                && current_page.endsWith('/index.html')) {
                link.classList.add('active');
                let parent = link.parentElement;
                while (parent) {
                    if (parent.tagName === 'LI' && parent.classList.contains('chapter-item')) {
                        parent.classList.add('expanded');
                    }
                    parent = parent.parentElement;
                }
            }
        }
        // Track and set sidebar scroll position
        this.addEventListener('click', e => {
            if (e.target.tagName === 'A') {
                const clientRect = e.target.getBoundingClientRect();
                const sidebarRect = this.getBoundingClientRect();
                sessionStorage.setItem('sidebar-scroll-offset', clientRect.top - sidebarRect.top);
            }
        }, { passive: true });
        const sidebarScrollOffset = sessionStorage.getItem('sidebar-scroll-offset');
        sessionStorage.removeItem('sidebar-scroll-offset');
        if (sidebarScrollOffset !== null) {
            // preserve sidebar scroll position when navigating via links within sidebar
            const activeSection = this.querySelector('.active');
            if (activeSection) {
                const clientRect = activeSection.getBoundingClientRect();
                const sidebarRect = this.getBoundingClientRect();
                const currentOffset = clientRect.top - sidebarRect.top;
                this.scrollTop += currentOffset - parseFloat(sidebarScrollOffset);
            }
        } else {
            // scroll sidebar to current active section when navigating via
            // 'next/previous chapter' buttons
            const activeSection = document.querySelector('#mdbook-sidebar .active');
            if (activeSection) {
                activeSection.scrollIntoView({ block: 'center' });
            }
        }
        // Toggle buttons
        const sidebarAnchorToggles = document.querySelectorAll('.chapter-fold-toggle');
        function toggleSection(ev) {
            ev.currentTarget.parentElement.parentElement.classList.toggle('expanded');
        }
        Array.from(sidebarAnchorToggles).forEach(el => {
            el.addEventListener('click', toggleSection);
        });
    }
}
window.customElements.define('mdbook-sidebar-scrollbox', MDBookSidebarScrollbox);


// ---------------------------------------------------------------------------
// Support for dynamically adding headers to the sidebar.

(function() {
    // This is used to detect which direction the page has scrolled since the
    // last scroll event.
    let lastKnownScrollPosition = 0;
    // This is the threshold in px from the top of the screen where it will
    // consider a header the "current" header when scrolling down.
    const defaultDownThreshold = 150;
    // Same as defaultDownThreshold, except when scrolling up.
    const defaultUpThreshold = 300;
    // The threshold is a virtual horizontal line on the screen where it
    // considers the "current" header to be above the line. The threshold is
    // modified dynamically to handle headers that are near the bottom of the
    // screen, and to slightly offset the behavior when scrolling up vs down.
    let threshold = defaultDownThreshold;
    // This is used to disable updates while scrolling. This is needed when
    // clicking the header in the sidebar, which triggers a scroll event. It
    // is somewhat finicky to detect when the scroll has finished, so this
    // uses a relatively dumb system of disabling scroll updates for a short
    // time after the click.
    let disableScroll = false;
    // Array of header elements on the page.
    let headers;
    // Array of li elements that are initially collapsed headers in the sidebar.
    // I'm not sure why eslint seems to have a false positive here.
    // eslint-disable-next-line prefer-const
    let headerToggles = [];
    // This is a debugging tool for the threshold which you can enable in the console.
    let thresholdDebug = false;

    // Updates the threshold based on the scroll position.
    function updateThreshold() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const windowHeight = window.innerHeight;
        const documentHeight = document.documentElement.scrollHeight;

        // The number of pixels below the viewport, at most documentHeight.
        // This is used to push the threshold down to the bottom of the page
        // as the user scrolls towards the bottom.
        const pixelsBelow = Math.max(0, documentHeight - (scrollTop + windowHeight));
        // The number of pixels above the viewport, at least defaultDownThreshold.
        // Similar to pixelsBelow, this is used to push the threshold back towards
        // the top when reaching the top of the page.
        const pixelsAbove = Math.max(0, defaultDownThreshold - scrollTop);
        // How much the threshold should be offset once it gets close to the
        // bottom of the page.
        const bottomAdd = Math.max(0, windowHeight - pixelsBelow - defaultDownThreshold);
        let adjustedBottomAdd = bottomAdd;

        // Adjusts bottomAdd for a small document. The calculation above
        // assumes the document is at least twice the windowheight in size. If
        // it is less than that, then bottomAdd needs to be shrunk
        // proportional to the difference in size.
        if (documentHeight < windowHeight * 2) {
            const maxPixelsBelow = documentHeight - windowHeight;
            const t = 1 - pixelsBelow / Math.max(1, maxPixelsBelow);
            const clamp = Math.max(0, Math.min(1, t));
            adjustedBottomAdd *= clamp;
        }

        let scrollingDown = true;
        if (scrollTop < lastKnownScrollPosition) {
            scrollingDown = false;
        }

        if (scrollingDown) {
            // When scrolling down, move the threshold up towards the default
            // downwards threshold position. If near the bottom of the page,
            // adjustedBottomAdd will offset the threshold towards the bottom
            // of the page.
            const amountScrolledDown = scrollTop - lastKnownScrollPosition;
            const adjustedDefault = defaultDownThreshold + adjustedBottomAdd;
            threshold = Math.max(adjustedDefault, threshold - amountScrolledDown);
        } else {
            // When scrolling up, move the threshold down towards the default
            // upwards threshold position. If near the bottom of the page,
            // quickly transition the threshold back up where it normally
            // belongs.
            const amountScrolledUp = lastKnownScrollPosition - scrollTop;
            const adjustedDefault = defaultUpThreshold - pixelsAbove
                + Math.max(0, adjustedBottomAdd - defaultDownThreshold);
            threshold = Math.min(adjustedDefault, threshold + amountScrolledUp);
        }

        if (documentHeight <= windowHeight) {
            threshold = 0;
        }

        if (thresholdDebug) {
            const id = 'mdbook-threshold-debug-data';
            let data = document.getElementById(id);
            if (data === null) {
                data = document.createElement('div');
                data.id = id;
                data.style.cssText = `
                    position: fixed;
                    top: 50px;
                    right: 10px;
                    background-color: 0xeeeeee;
                    z-index: 9999;
                    pointer-events: none;
                `;
                document.body.appendChild(data);
            }
            data.innerHTML = `
                <table>
                  <tr><td>documentHeight</td><td>${documentHeight.toFixed(1)}</td></tr>
                  <tr><td>windowHeight</td><td>${windowHeight.toFixed(1)}</td></tr>
                  <tr><td>scrollTop</td><td>${scrollTop.toFixed(1)}</td></tr>
                  <tr><td>pixelsAbove</td><td>${pixelsAbove.toFixed(1)}</td></tr>
                  <tr><td>pixelsBelow</td><td>${pixelsBelow.toFixed(1)}</td></tr>
                  <tr><td>bottomAdd</td><td>${bottomAdd.toFixed(1)}</td></tr>
                  <tr><td>adjustedBottomAdd</td><td>${adjustedBottomAdd.toFixed(1)}</td></tr>
                  <tr><td>scrollingDown</td><td>${scrollingDown}</td></tr>
                  <tr><td>threshold</td><td>${threshold.toFixed(1)}</td></tr>
                </table>
            `;
            drawDebugLine();
        }

        lastKnownScrollPosition = scrollTop;
    }

    function drawDebugLine() {
        if (!document.body) {
            return;
        }
        const id = 'mdbook-threshold-debug-line';
        const existingLine = document.getElementById(id);
        if (existingLine) {
            existingLine.remove();
        }
        const line = document.createElement('div');
        line.id = id;
        line.style.cssText = `
            position: fixed;
            top: ${threshold}px;
            left: 0;
            width: 100vw;
            height: 2px;
            background-color: red;
            z-index: 9999;
            pointer-events: none;
        `;
        document.body.appendChild(line);
    }

    function mdbookEnableThresholdDebug() {
        thresholdDebug = true;
        updateThreshold();
        drawDebugLine();
    }

    window.mdbookEnableThresholdDebug = mdbookEnableThresholdDebug;

    // Updates which headers in the sidebar should be expanded. If the current
    // header is inside a collapsed group, then it, and all its parents should
    // be expanded.
    function updateHeaderExpanded(currentA) {
        // Add expanded to all header-item li ancestors.
        let current = currentA.parentElement;
        while (current) {
            if (current.tagName === 'LI' && current.classList.contains('header-item')) {
                current.classList.add('expanded');
            }
            current = current.parentElement;
        }
    }

    // Updates which header is marked as the "current" header in the sidebar.
    // This is done with a virtual Y threshold, where headers at or below
    // that line will be considered the current one.
    function updateCurrentHeader() {
        if (!headers || !headers.length) {
            return;
        }

        // Reset the classes, which will be rebuilt below.
        const els = document.getElementsByClassName('current-header');
        for (const el of els) {
            el.classList.remove('current-header');
        }
        for (const toggle of headerToggles) {
            toggle.classList.remove('expanded');
        }

        // Find the last header that is above the threshold.
        let lastHeader = null;
        for (const header of headers) {
            const rect = header.getBoundingClientRect();
            if (rect.top <= threshold) {
                lastHeader = header;
            } else {
                break;
            }
        }
        if (lastHeader === null) {
            lastHeader = headers[0];
            const rect = lastHeader.getBoundingClientRect();
            const windowHeight = window.innerHeight;
            if (rect.top >= windowHeight) {
                return;
            }
        }

        // Get the anchor in the summary.
        const href = '#' + lastHeader.id;
        const a = [...document.querySelectorAll('.header-in-summary')]
            .find(element => element.getAttribute('href') === href);
        if (!a) {
            return;
        }

        a.classList.add('current-header');

        updateHeaderExpanded(a);
    }

    // Updates which header is "current" based on the threshold line.
    function reloadCurrentHeader() {
        if (disableScroll) {
            return;
        }
        updateThreshold();
        updateCurrentHeader();
    }


    // When clicking on a header in the sidebar, this adjusts the threshold so
    // that it is located next to the header. This is so that header becomes
    // "current".
    function headerThresholdClick(event) {
        // See disableScroll description why this is done.
        disableScroll = true;
        setTimeout(() => {
            disableScroll = false;
        }, 100);
        // requestAnimationFrame is used to delay the update of the "current"
        // header until after the scroll is done, and the header is in the new
        // position.
        requestAnimationFrame(() => {
            requestAnimationFrame(() => {
                // Closest is needed because if it has child elements like <code>.
                const a = event.target.closest('a');
                const href = a.getAttribute('href');
                const targetId = href.substring(1);
                const targetElement = document.getElementById(targetId);
                if (targetElement) {
                    threshold = targetElement.getBoundingClientRect().bottom;
                    updateCurrentHeader();
                }
            });
        });
    }

    // Takes the nodes from the given head and copies them over to the
    // destination, along with some filtering.
    function filterHeader(source, dest) {
        const clone = source.cloneNode(true);
        clone.querySelectorAll('mark').forEach(mark => {
            mark.replaceWith(...mark.childNodes);
        });
        dest.append(...clone.childNodes);
    }

    // Scans page for headers and adds them to the sidebar.
    document.addEventListener('DOMContentLoaded', function() {
        const activeSection = document.querySelector('#mdbook-sidebar .active');
        if (activeSection === null) {
            return;
        }

        const main = document.getElementsByTagName('main')[0];
        headers = Array.from(main.querySelectorAll('h2, h3, h4, h5, h6'))
            .filter(h => h.id !== '' && h.children.length && h.children[0].tagName === 'A');

        if (headers.length === 0) {
            return;
        }

        // Build a tree of headers in the sidebar.

        const stack = [];

        const firstLevel = parseInt(headers[0].tagName.charAt(1));
        for (let i = 1; i < firstLevel; i++) {
            const ol = document.createElement('ol');
            ol.classList.add('section');
            if (stack.length > 0) {
                stack[stack.length - 1].ol.appendChild(ol);
            }
            stack.push({level: i + 1, ol: ol});
        }

        // The level where it will start folding deeply nested headers.
        const foldLevel = 3;

        for (let i = 0; i < headers.length; i++) {
            const header = headers[i];
            const level = parseInt(header.tagName.charAt(1));

            const currentLevel = stack[stack.length - 1].level;
            if (level > currentLevel) {
                // Begin nesting to this level.
                for (let nextLevel = currentLevel + 1; nextLevel <= level; nextLevel++) {
                    const ol = document.createElement('ol');
                    ol.classList.add('section');
                    const last = stack[stack.length - 1];
                    const lastChild = last.ol.lastChild;
                    // Handle the case where jumping more than one nesting
                    // level, which doesn't have a list item to place this new
                    // list inside of.
                    if (lastChild) {
                        lastChild.appendChild(ol);
                    } else {
                        last.ol.appendChild(ol);
                    }
                    stack.push({level: nextLevel, ol: ol});
                }
            } else if (level < currentLevel) {
                while (stack.length > 1 && stack[stack.length - 1].level > level) {
                    stack.pop();
                }
            }

            const li = document.createElement('li');
            li.classList.add('header-item');
            li.classList.add('expanded');
            if (level < foldLevel) {
                li.classList.add('expanded');
            }
            const span = document.createElement('span');
            span.classList.add('chapter-link-wrapper');
            const a = document.createElement('a');
            span.appendChild(a);
            a.href = '#' + header.id;
            a.classList.add('header-in-summary');
            filterHeader(header.children[0], a);
            a.addEventListener('click', headerThresholdClick);
            const nextHeader = headers[i + 1];
            if (nextHeader !== undefined) {
                const nextLevel = parseInt(nextHeader.tagName.charAt(1));
                if (nextLevel > level && level >= foldLevel) {
                    const toggle = document.createElement('a');
                    toggle.classList.add('chapter-fold-toggle');
                    toggle.classList.add('header-toggle');
                    toggle.addEventListener('click', () => {
                        li.classList.toggle('expanded');
                    });
                    const toggleDiv = document.createElement('div');
                    toggleDiv.textContent = '❱';
                    toggle.appendChild(toggleDiv);
                    span.appendChild(toggle);
                    headerToggles.push(li);
                }
            }
            li.appendChild(span);

            const currentParent = stack[stack.length - 1];
            currentParent.ol.appendChild(li);
        }

        const onThisPage = document.createElement('div');
        onThisPage.classList.add('on-this-page');
        onThisPage.append(stack[0].ol);
        const activeItemSpan = activeSection.parentElement;
        activeItemSpan.after(onThisPage);
    });

    document.addEventListener('DOMContentLoaded', reloadCurrentHeader);
    document.addEventListener('scroll', reloadCurrentHeader, { passive: true });
})();

