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

  function setStyles(element, styles) {
    Object.keys(styles).forEach(function (key) {
      element.style[key] = styles[key];
    });
    return element;
  }

  function charColor(ch) {
    if (ch === "_") return "#737373";
    if (ch === "o" || ch === "O" || ch === "0") return "#facc15";
    return "#f5f5f5";
  }

  function renderLogo() {
    var index = 0;
    var fragment = document.createDocumentFragment();

    logo.forEach(function (line, row) {
      Array.prototype.forEach.call(line, function (ch, column) {
        var span = document.createElement("span");
        span.textContent = ch;
        setStyles(span, {
          color: charColor(ch),
          display: "inline-block",
          opacity: "0",
          textShadow: ch === "o" || ch === "O" || ch === "0" ? "0 0 16px rgb(250 204 21 / 0.42)" : "none",
          transform: "translateY(8px) scale(0.98)"
        });
        span.dataset.charIndex = String(index);
        span.dataset.charColumn = String(column);
        span.dataset.charRow = String(row);
        fragment.appendChild(span);
        index += 1;
      });
      if (row < logo.length - 1) {
        fragment.appendChild(document.createTextNode("\n"));
      }
    });

    return fragment;
  }

  function removeIntro(intro, previousOverflow) {
    intro.style.pointerEvents = "none";
    intro.animate([{ opacity: 1 }, { opacity: 0 }], {
      duration: 360,
      easing: "ease",
      fill: "forwards"
    });
    document.body.style.overflow = previousOverflow;
    window.setTimeout(function () {
      intro.remove();
    }, 460);
  }

  function mountIntro() {
    var previousOverflow = document.body.style.overflow;
    var intro = setStyles(document.createElement("div"), {
      background: "linear-gradient(180deg, rgb(255 255 255 / 0.04), transparent 42%), #050505",
      color: "#f5f5f5",
      display: "grid",
      inset: "0",
      opacity: "1",
      placeItems: "center",
      position: "fixed",
      zIndex: "100"
    });
    intro.className = "docs-intro";
    intro.setAttribute("role", "dialog");
    intro.setAttribute("aria-label", "Carbide opening animation");

    var stage = setStyles(document.createElement("div"), {
      transform: "translateY(0)",
      width: "min(920px, calc(100vw - 32px))"
    });

    var logoNode = setStyles(document.createElement("pre"), {
      background: "transparent",
      border: "0",
      color: "#f5f5f5",
      fontSize: "clamp(8px, 1.45vw, 15px)",
      letterSpacing: "0",
      lineHeight: "1.08",
      margin: "0 auto",
      maxWidth: "100%",
      overflow: "visible",
      padding: "0",
      textAlign: "left",
      textShadow: "0 20px 50px rgb(0 0 0 / 0.8)",
      whiteSpace: "pre",
      width: "max-content"
    });
    logoNode.setAttribute("aria-label", "Carbide ASCII logo");
    logoNode.appendChild(renderLogo());

    stage.appendChild(logoNode);
    intro.appendChild(stage);

    document.body.prepend(intro);
    document.body.style.overflow = "hidden";

    stage.animate([{ opacity: 0, transform: "translateY(10px) scale(0.985)" }, { opacity: 1, transform: "translateY(0) scale(1)" }], {
      duration: 520,
      easing: "cubic-bezier(0.2, 0.8, 0.2, 1)",
      fill: "both"
    });
    logoNode.querySelectorAll("span").forEach(function (span) {
      var delay = 160 + (Number(span.dataset.charRow) * 48) + (Number(span.dataset.charColumn) * 4);
      span.animate([{ opacity: 0, transform: "translateY(8px) scale(0.98)" }, { opacity: 1, transform: "translateY(0) scale(1)" }], {
        delay: delay,
        duration: 440,
        easing: "cubic-bezier(0.2, 0.9, 0.2, 1)",
        fill: "both"
      });
    });

    window.setTimeout(function () {
      removeIntro(intro, previousOverflow);
    }, 1850);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", mountIntro, { once: true });
  } else {
    mountIntro();
  }
}());
