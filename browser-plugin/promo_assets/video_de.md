---
marp: true
theme: default
size: 16:9
style: |
  /* --- LAYOUT CONTAINER --- */
  section.small-layout {
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
    position: relative;
    overflow: hidden;
    font-family: 'Helvetica', 'Arial', sans-serif;
    color: #ffffff;
  }

  /* --- TEXT STYLING --- */
  .text-content {
    position: absolute;
    left: 60px;
    top: 25%;
    width: 35%;
    z-index: 10;
  }

  .text-content h1 {
    font-size: 60px;
    color: #ffffff;
    margin-bottom: 20px;
    line-height: 1.1;
  }

  .text-content p {
    font-size: 32px;
    color: #ffffff;
    line-height: 1.4;
  }

  /* --- MAIN BROWSER IMAGE --- */
  img.main-browser {
    position: absolute;
    right: -20px;
    top: 80px;
    width: 58%;
    box-shadow: -10px 10px 30px rgba(0,0,0,0.15);
    border-radius: 12px;
    border: 1px solid rgba(0,0,0,0.1);
    z-index: 1;
  }

  /* --- SMALL OVERLAY IMAGE --- */
  img.small-overlay {
    position: absolute;
    bottom: 50px;
    left: 34%;
    width: 34%;
    /* HEAVY SHADOW ADDED BELOW */
    box-shadow: 0 35px 70px rgba(0,0,0,0.7);
    border-radius: 12px;
    border: 3px solid #ffffff;
    z-index: 5;
  }

  /* --- LOGO --- */
  img.logo {
    position: absolute;
    bottom: 40px;
    right: 40px;
    height: 80px;
    width: auto;
    z-index: 10;
  }

---

<!-- _class: small-layout -->

<div class="text-content">

# News Deframer

Erkenne und vermeide Clickbait auf deinen Nachrichtenportalen

</div>

<img src="assets/screenshot2.png" class="main-browser" />
<img src="assets/screenshot1.png" class="small-overlay" />
<img src="assets/logo.svg" class="logo" />
