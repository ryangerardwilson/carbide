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

    var track = setStyles(document.createElement("div"), {
      height: "18px",
      margin: "0 0 18px",
      overflow: "hidden",
      position: "relative"
    });
    track.setAttribute("aria-hidden", "true");

    var pellets = setStyles(document.createElement("span"), {
      backgroundImage: "radial-gradient(circle, rgb(250 204 21 / 0.7) 0 2px, transparent 2.5px)",
      backgroundSize: "18px 2px",
      height: "2px",
      inset: "8px 0 auto",
      opacity: "0.72",
      position: "absolute"
    });

    var chomper = setStyles(document.createElement("span"), {
      background: "#facc15",
      borderRadius: "999px",
      boxShadow: "0 0 18px rgb(250 204 21 / 0.48)",
      clipPath: "polygon(0 0, 100% 0, 58% 50%, 100% 100%, 0 100%)",
      height: "18px",
      left: "0",
      position: "absolute",
      top: "0",
      width: "18px"
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

    var footer = setStyles(document.createElement("div"), {
      alignItems: "center",
      color: "#a3a3a3",
      display: "flex",
      fontSize: "0.78rem",
      fontWeight: "700",
      gap: "16px",
      justifyContent: "center",
      marginTop: "24px",
      textTransform: "uppercase"
    });
    var footerText = document.createElement("span");
    footerText.textContent = "Carbide docs";
    footer.appendChild(footerText);

    track.appendChild(pellets);
    track.appendChild(chomper);
    stage.appendChild(track);
    stage.appendChild(logoNode);
    stage.appendChild(footer);
    intro.appendChild(stage);

    document.body.prepend(intro);
    document.body.style.overflow = "hidden";

    stage.animate([{ opacity: 0, transform: "translateY(10px)" }, { opacity: 1, transform: "translateY(0)" }], {
      duration: 520,
      easing: "cubic-bezier(0.2, 0.8, 0.2, 1)",
      fill: "both"
    });
    chomper.animate([{ transform: "translateX(-22px)" }, { transform: "translateX(calc(min(920px, calc(100vw - 32px)) - 18px))" }], {
      duration: 1320,
      easing: "cubic-bezier(0.55, 0, 0.1, 1)",
      fill: "both"
    });
    chomper.animate([
      { clipPath: "polygon(0 0, 100% 0, 58% 50%, 100% 100%, 0 100%)" },
      { clipPath: "circle(50% at 50% 50%)" },
      { clipPath: "polygon(0 0, 100% 0, 58% 50%, 100% 100%, 0 100%)" }
    ], {
      duration: 180,
      easing: "steps(2, end)",
      iterations: Infinity
    });
    logoNode.querySelectorAll("span").forEach(function (span) {
      span.animate([{ opacity: 0, transform: "translateY(8px) scale(0.98)" }, { opacity: 1, transform: "translateY(0) scale(1)" }], {
        delay: 220 + (Number(span.dataset.charIndex) * 0.55),
        duration: 440,
        easing: "cubic-bezier(0.2, 0.9, 0.2, 1)",
        fill: "both"
      });
    });

    window.setTimeout(function () {
      removeIntro(intro, previousOverflow);
    }, 2450);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", mountIntro, { once: true });
  } else {
    mountIntro();
  }
}());
