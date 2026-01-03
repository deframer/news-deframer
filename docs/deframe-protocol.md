# The Deframe Protocol
### A Decentralized Cognitive Security Layer on the AT Protocol

**Version:** 0.0 (Vision)
**Target Network:** Bluesky / AT Protocol

---

## 1. Abstract
The Deframe Protocol is a distributed system designed to analyze and annotate information quality on the social web. Unlike traditional fact-checking, which seeks a binary "True/False" verdict, Deframe analyzes the *nature* of the content—measuring cognitive stressors such as **Clickbait**, **Persuasive Intent**, **Hyper-Stimulus**, and **Speculative Content**.

The system utilizes the Bluesky (AT Protocol) network as a public, immutable ledger. It employs a swarm of independent worker bots to generate consensus ratings, calculates reputation over time, and mathematically isolates malicious or "hallucinating" actors.

---

## 2. Core Architecture
The system consists of four distinct actors operating in a public Bluesky thread.

### 2.1 The Actors

1.  **The Trusted Flagger (Source)**
    *   **Role:** Identifies content that needs analysis.
    *   **Action:** Canonicalizes the URL and data, calculates a SHA256 checksum, and posts the "Job" to the network.
2.  **The Control Bot (Auditor)**
    *   **Role:** The Gatekeeper.
    *   **Action:** Immediately validates that the Job's payload matches its SHA256 ID. If valid, it posts a "Green Light" reply. This prevents the swarm from processing corrupted or spam data.
3.  **The Worker Swarm (Raters)**
    *   **Role:** The Analysts.
    *   **Action:** Independent agents that pick up valid jobs, perform analysis (via AI, heuristics, or human input), and post their ratings as a reply.
4.  **The Consensus Engine (Aggregator)**
    *   **Role:** The Judge.
    *   **Action:** Listens for worker replies, calculates the weighted average, updates worker reputation scores based on their accuracy, and publishes the final result.

---

## 3. The Cognitive Attributes
Workers do not determine "Truth." They measure specific psychological and rhetorical attributes on a scale of `0.0` (None) to `1.0` (Extreme).

*   **Clickbait:** The gap between the headline's promise and the content's delivery.
*   **Persuasive Intent:** The degree to which the text attempts to change behavior or opinion rather than inform.
*   **Hyper-Stimulus:** The density of emotional triggers, capitalization, urgency, and inflammatory language.
*   **Speculative Content:** The ratio of verifiable facts to unverified claims or future predictions.

---

## 4. Technical Implementation (AT Protocol)

This system leverages specific features of Bluesky to ensure data integrity and uniqueness.

### 4.1 The Unique Job ID (The RKey)
To prevent duplicate jobs and ensure a universal reference, the **SHA256 hash** of the canonical payload is used as the **Record Key (`rkey`)** in the AT Protocol.

*   **Repository:** `@flagger.deframe.net`
*   **Collection:** `app.bsky.feed.post`
*   **RKey:** `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` (The SHA256)

If a Flagger tries to post the same URL/Job twice, the protocol rejects it automatically.

### 4.2 The Data Flow (ASCII Diagram)

```text
[ TRUSTED FLAGGER ]
       |
       | 1. Create Job (SHA256 = RKey)
       v
 +------------------------------------------+
 | ROOT POST (ID: sha256...)                |
 | "Analyze this content"                   |
 +------------------------------------------+
       |
       | 2. Watches Stream
       v
 [ CONTROL BOT ]
       | Checks Integrity... OK!
       |
       +---> [ REPLY 1 ] "✅ VALIDATED"
                  |
                  v
           [ WORKER SWARM ]
           (Workers wait for "VALIDATED" signal)
           (Each performs independent analysis)
                  |
        +---------+---------+
        |         |         |
   [ REPLY 2 ] [ REPLY 3 ] [ REPLY 4 ]
   {Ratings}   {Ratings}   {Ratings}
        |         |         |
        v         v         v
      [ CONSENSUS ENGINE ]
      (Wait for X replies or T time)
      (Calculate Weighted Averages)
      (Update Reputation DB)
                  |
                  v
          [ FINAL POST ]
          "🏁 CONSENSUS REACHED"
          Clickbait: 0.8 | Persuasive: 0.9
```

---

## 5. The Reputation Model (Learning)

The system does not require code changes to detect bad bots. It uses a mathematical approach to "Learn" trust over time.

### 5.1 The Logic: "Distance from Consensus"
The "Truth" is defined as the weighted average of the swarm.
1.  Calculate the **Average** of all ratings for a job.
2.  Calculate the **Deviation** (Distance) of each bot from that average.
3.  **Reward** bots close to the average.
4.  **Penalize** bots far from the average.

### 5.2 Example Scenario
*   **Job:** A clear scam article.
*   **Bot A:** Rates Risk `0.9` (Correct)
*   **Bot B:** Rates Risk `0.85` (Correct)
*   **Bot C (Malicious):** Rates Risk `0.1` (Trying to hide the scam)

**The Math:**
*   **Average:** `(0.9 + 0.85 + 0.1) / 3 = 0.61`
*   **Bot A Deviation:** `|0.9 - 0.61| = 0.29` (Acceptable)
*   **Bot C Deviation:** `|0.1 - 0.61| = 0.51` (High Error)

**The Consequence:**
Bot C loses "Trust Points." In the next job, Bot C's vote will carry less weight (e.g., only 10% influence), effectively silencing malicious actors without banning them.

---

## 6. JSON Schemas

### 6.1 The Job Payload (Flagger)
*Embedded in the Root Post.*
```json
{
  "type": "deframe_job",
  "target_url": "https://example-news.com/article",
  "sha256": "e3b0c44...",
  "instructions": "standard_analysis"
}
```

### 6.2 The Worker Rating (Reply)
*Posted by Worker Bots.*
```json
{
  "type": "deframe_rating",
  "ref_root": "at://did:plc:123.../app.bsky.feed.post/e3b0c44...",
  "attributes": {
    "clickbait": 0.8,
    "persuasive_intent": 0.5,
    "hyper_stimulus": 0.2,
    "speculative_content": 0.9
  }
}
```

### 6.3 The Consensus Result (Final)
*Posted by the Consensus Engine.*
```json
{
  "type": "deframe_consensus",
  "participant_count": 12,
  "trust_score": 0.98, // Confidence in the result
  "final_attributes": {
    "clickbait": 0.79,
    "persuasive_intent": 0.51,
    "hyper_stimulus": 0.22,
    "speculative_content": 0.88
  }
}
```

---

## 7. Conclusion
By decoupling the **Source** (Flagger), the **Analysis** (Workers), and the **Verdict** (Consensus), the Deframe Protocol creates a resilient, censorship-resistant layer for understanding content quality. It turns the "noise" of social media into structured, verifiable data, leveraging the transparency of the Bluesky network.
