---
marp: true
theme: default
size: 16:9
style: |
  /* --- LAYOUT CONTAINER --- */
  section.small-layout {
    /* Dark blue gradient background */
    background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
    position: relative;
    overflow: hidden; /* Cuts off images that slide off the edge */
    font-family: 'Helvetica', 'Arial', sans-serif;
    color: #ffffff;
  }

  /* --- TEXT STYLING --- */
  .text-content {
    position: absolute;
    left: 60px;
    top: 25%;
    width: 35%; /* Restrict text to the left side */
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

  /* --- MAIN BROWSER IMAGE (BACKGROUND) --- */
  img.main-browser {
    position: absolute;
    right: -100px; /* Push slightly off screen */
    top: 80px;
    width: 65%;
    box-shadow: -10px 10px 30px rgba(0,0,0,0.15);
    border-radius: 12px;
    border: 1px solid rgba(0,0,0,0.1);
    z-index: 1; /* Sits behind */
  }

  /* --- SMALL OVERLAY IMAGE (FOREGROUND) --- */
  img.small-overlay {
    position: absolute;
    bottom: 50px;
    left: 30%; /* Positions it bridging the text and browser */
    width: 40%;

    /* Styling to make it pop */
    box-shadow: 0 20px 40px rgba(0,0,0,0.4);
    border-radius: 12px;
    border: 3px solid #ffffff; /* White border like the reference */
    z-index: 5; /* Sits on top */
  }

  /* Style the Logo in bottom right */
  img.logo {
    position: absolute;
    bottom: 30px;
    right: 40px;
    height: 100px; /* Adjust based on your logo aspect ratio */
    width: auto;
  }

---

<!-- _class: small-layout -->

<div class="text-content">

# News Deframer

Detect and avoid clickbait from your favorite news portal

</div>

<img src="assets/screenshot2.png" class="main-browser" />

<img src="assets/screenshot1.png" class="small-overlay" />

<img src="assets/logo.svg" class="logo" />