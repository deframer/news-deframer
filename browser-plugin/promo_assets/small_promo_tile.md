---
marp: true
theme: default
size: 16:10
style: |
  /* --- LAYOUT CONTAINER --- */
  section.small-layout {
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
    position: relative;
    overflow: hidden;
    font-family: 'Helvetica', 'Arial', sans-serif;
    color: #ffffff;
    padding: 0;
  }

  /* --- TEXT BOX --- */
  .text-box {
    position: absolute;
    left: 80px;
    top: 40%;
    transform: translateY(-50%);
    width: 42%;

    /* Glass Effect */
    background: rgba(15, 23, 42, 0.85);
    backdrop-filter: blur(8px);

    padding: 40px;
    border-radius: 16px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    box-shadow: 0 20px 40px rgba(0,0,0,0.4);
    z-index: 20; /* High Z-Index ensures text sits ON TOP of the large images */
  }

  .text-box h1 {
    font-size: 48px;
    font-weight: bold;
    color: #ffffff;
    margin: 0 0 20px 0;
    line-height: 1.1;
  }

  .text-box p {
    font-size: 26px;
    color: #cbd5e1;
    margin: 0;
    line-height: 1.4;
  }

  /* --- MAIN BROWSER IMAGE --- */
  img.main-browser {
    position: absolute;
    right: -20px; /* Moved Left (closer to center) */
    top: 80px;
    width: 75%;   /* MUCH Larger (was 55%) */

    box-shadow: -10px 10px 30px rgba(0,0,0,0.3);
    border-radius: 12px;
    opacity: 0.9;
    z-index: 1;   /* Sits behind text box */
  }

  /* --- SMALL OVERLAY IMAGE --- */
  img.small-overlay {
    position: absolute;
    bottom: 40px;
    left: 40%;    /* Moved Left (was 48%) */
    width: 42%;   /* Larger (was 30%) */

    box-shadow: 0 20px 40px rgba(0,0,0,0.5);
    border-radius: 12px;
    border: 3px solid #ffffff;
    z-index: 5;
  }

  /* --- LOGO --- */
  img.logo {
    position: absolute;
    bottom: 40px;
    right: 100px;
    height: 70px;
    width: auto;
    z-index: 10;
  }
---

<!-- _class: small-layout -->

<div class="text-box">

# News Deframer

Detect and avoid clickbait from your favorite news portal

</div>

<img src="assets/screenshot2.png" class="main-browser" />
<img src="assets/screenshot1.png" class="small-overlay" />
<img src="assets/logo.svg" class="logo" />
