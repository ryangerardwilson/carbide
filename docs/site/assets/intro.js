(function () {
  "use strict";

  var replayRequested = window.location.search.indexOf("intro=1") !== -1;
  var reducedMotion = window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  function isHomePage() {
    var pathname = window.location.pathname;
    return pathname === "/" || pathname === "/index.html";
  }

  if (reducedMotion || (!replayRequested && !isHomePage())) {
    return;
  }

  var logo = [
    "_____________________________________________________",
    "________________________oo_______oo_______oo_________",
    "_ooooo___ooooo__oo_ooo__oooooo________oooooo__ooooo__",
    "oo___oo_oo___oo_ooo___o_oo___oo__oo__oo___oo_oo____o_",
    "oo______oo___oo_oo______oo___oo__oo__oo___oo_ooooooo_",
    "oo______oo___oo_oo______oo___oo__oo__oo___oo_oo______",
    "_ooooo___oooo_o_oo______oooooo__oooo__oooooo__ooooo__",
    "_____________________________________________________"
  ];

  function colorClass(ch) {
    if (ch === "_") {
      return "logo-char logo-char-rail";
    }
    if (ch === "o" || ch === "O" || ch === "0") {
      return "logo-char logo-char-round";
    }
    return "logo-char";
  }

  function renderLogo() {
    var index = 0;
    return logo.map(function (line, row) {
      var chars = Array.prototype.map.call(line, function (ch, column) {
        var span = document.createElement("span");
        span.className = colorClass(ch);
        span.style.setProperty("--char-index", String(index));
        span.style.setProperty("--char-column", String(column));
        span.style.setProperty("--char-row", String(row));
        span.textContent = ch;
        index += 1;
        return span.outerHTML;
      }).join("");
      return chars;
    }).join("\n");
  }

  function removeIntro(intro) {
    intro.classList.add("docs-intro-exit");
    document.body.classList.remove("intro-active");
    window.setTimeout(function () {
      intro.remove();
    }, 460);
  }

  function mountIntro() {
    var intro = document.createElement("div");
    intro.className = "docs-intro";
    intro.setAttribute("role", "dialog");
    intro.setAttribute("aria-label", "Carbide opening animation");
    intro.innerHTML = [
      '<div class="docs-intro-stage">',
      '  <div class="docs-intro-track" aria-hidden="true">',
      '    <span class="docs-intro-chomper"></span>',
      '  </div>',
      '  <pre class="docs-intro-logo" aria-label="Carbide ASCII logo">' + renderLogo() + "</pre>",
      '  <div class="docs-intro-footer">',
      '    <span>Carbide docs</span>',
      "  </div>",
      "</div>"
    ].join("");

    document.body.prepend(intro);
    document.body.classList.add("intro-active");

    window.setTimeout(function () {
      removeIntro(intro);
    }, 2450);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", mountIntro, { once: true });
  } else {
    mountIntro();
  }
}());
