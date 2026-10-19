COURSE = dict(
    key="c2",
    title="Investing on the ESX: A Practical Guide",
    duration="2 hours 30 mins",
    description="A hands-on guide to the Ethiopian Securities Exchange: the seven companies listed so far, how to open an account and place your first order, how prices and order books work, what trading really costs, how IPOs and prospectuses work, what dividends and AGMs mean for you, and how to avoid scams.",
)

LESSONS = [
dict(title="The ESX at a Glance", video="", content="""
## What the ESX is

The **Ethiopian Securities Exchange (ESX)** is Ethiopia's stock exchange. It is a licensed marketplace where securities are listed, bought and sold under rules approved by ECMA. Its office is at the Nile Insurance S.C. Building (16th to 18th floors), National Theatre area, Addis Ababa.

## Listed companies (as of October 2026)

All seven are on the **Main Market** and are listed as equities:

| Ticker | Company | Sector | Listed on |
|---|---|---|---|
| WGBX | Wegagen Bank | Banking | 10 Jan 2025 |
| GDAB | Gadaa Bank | Banking | 13 Jun 2025 |
| AWAB | Awash Bank | Banking | 23 Apr 2026 |
| TELE | Ethio Telecom | Telecom | 26 May 2026 |
| ABAYB | Abay Bank | Banking | 25 Jun 2026 |
| BOAX | Bank of Abyssinia | Banking | 28 Jul 2026 |
| SIDAX | Sidama Bank | Banking | 28 Sep 2026 |

You can see each company's profile, board, audited financials and ESX filings on Stocet's **Fundamentals** pages.

## What to notice

- **It is concentrated.** Six of the seven are banks, so the market's performance is heavily tied to the banking sector and NBE policy.
- **It is young.** The first listing was less than two years ago. Track records, trading volumes and market-making are still developing.
- **It is growing.** Other companies (insurers, other banks, manufacturers) are preparing to list.
- **Ticker symbols** are short codes (like AWAB) used to place orders and to find a company on Stocet. Some third-party sites use different codes; the ESX code is the official one.

## Main Market vs. other segments

The **Main Market** is for larger, established companies meeting full listing requirements. ESX also provides for other segments (such as a growth market and fixed-income trading) as the exchange develops. Always check which segment a security trades in.

## Key takeaways

1. Seven companies are listed so far: six banks and Ethio Telecom.
2. The exchange is young, so liquidity and diversification are limited.
3. The ESX ticker is the official identifier for every listed company.
""", quiz=[
 ("Which company was the first to list on the ESX?", ["Awash Bank", "Ethio Telecom", "Wegagen Bank", "Gadaa Bank"], 2),
 ("How many of the seven listed companies are banks?", ["Three", "Four", "Six", "Seven"], 2),
 ("What does this concentration mean for an investor?", ["Returns are very diversified", "Market performance is closely tied to the banking sector", "Shares cannot be sold", "There is no regulator"], 1),
 ("What is the ESX ticker for Bank of Abyssinia?", ["BOA", "BOAX", "ABAYB", "BABY"], 1),
 ("Which company listed most recently (as of October 2026)?", ["Abay Bank", "Bank of Abyssinia", "Sidama Bank", "Ethio Telecom"], 2),
]),

dict(title="Opening an Account and Placing Your First Order", video="", content="""
## You cannot trade with the exchange directly

Retail investors place orders through a **licensed securities dealer / broker**. Today that means a licensed intermediary such as **Ethio-Fidelity Securities**, **Gadaa Securities Dealer**, or an investment bank that also offers brokerage such as **CBE Capital**. Check the up-to-date list on Stocet's CMSPs page and on ECMA's website.

## Step by step

1. **Choose a licensed firm.** Compare fees, how you place orders (app, phone, branch), customer support and minimum deposits. Confirm the firm is on ECMA's list.
2. **Complete KYC (Know Your Customer).** You will provide your national ID, contact details and bank account information, and sign the account agreement.
3. **Get your accounts.** You are given a trading account with the firm and a securities (depository) account that records the shares you own.
4. **Fund your account.** Transfer money from your bank to the dealer's designated account, following their instructions exactly.
5. **Choose a security** and research it (see the Fundamentals and Analytics sections of Stocet).
6. **Place an order** via the dealer's channel, which can include the **ESX Neway** trading app. Stocet itself never places orders: when you plan a trade in Stocet Analytics, the *Execute* button copies your plan and takes you to Neway.
7. **Confirm the trade** and check that the shares appear in your securities account.

## Order types you will meet

- **Market order**: buy or sell immediately at the best available price. Fast, but the price may differ from what you saw.
- **Limit order**: buy or sell **only at your price or better**. You control the price but the order may not execute.

Where available, use **limit orders**, especially in thinly traded shares, to avoid paying more than you intended.

## Before you click "buy"

Write down: *why* you are buying, the price you will not pay above, how much of your money this is, and at what price or situation you would sell. (The Risk Management course covers this.)

## Key takeaways

1. You trade through a licensed dealer, never directly with the ESX.
2. Expect KYC, a trading account, a securities account and a funding step.
3. Prefer limit orders in thin markets so you control your price.
4. Plan the exit before you enter.
""", quiz=[
 ("How do retail investors normally buy ESX-listed shares?", ["By walking into the ESX office", "Through a licensed securities dealer or broker", "Through the NBE", "By emailing the company"], 1),
 ("What does a limit order guarantee?", ["The order will always execute", "You will never pay more (or sell for less) than your set price", "A profit", "A lower fee"], 1),
 ("What does KYC stand for?", ["Keep Your Cash", "Know Your Customer", "Key Yield Calculation", "Know Your Company"], 1),
 ("Does Stocet place trades for you?", ["Yes, with one click", "No, it prepares your plan and you place the order with your licensed dealer / ESX Neway", "Only for Premium users", "Only for IPOs"], 1),
 ("Why are limit orders especially useful in thinly traded shares?", ["They remove all risk", "They stop you paying far above the last price when few sellers exist", "They are free", "They execute faster"], 1),
]),

dict(title="Prices, Order Books and Trading Sessions", video="", content="""
## How a price is formed

A price on an exchange is simply where a **buyer** and a **seller** agree. Every participant posts orders:

- **Bid**: the highest price a buyer will pay.
- **Ask (offer)**: the lowest price a seller will accept.
- **Spread**: the gap between the best bid and the best ask. A **narrow** spread usually means a liquid market; a **wide** spread means few participants and higher trading costs.

The **order book** lists all bids and asks with the quantities at each price. A trade happens when a bid and an ask meet.

*Example:* best bid ETB 1,150 and best ask ETB 1,180 means a spread of ETB 30 (about 2.6%). If you buy at the ask and immediately sell at the bid, you lose the spread.

## Trading sessions

According to the ESX's trading operations, the default trading day runs **9:00 to 15:00 East Africa Time, on weekdays**. A typical day has a **pre-open** session (about 9:00 to 9:30, when orders are collected and an opening price is determined) followed by **continuous trading** (9:30 to 15:00). ESX may adjust hours by notice, so check the exchange's announcements.

## Price movement limits

Exchanges limit how far a price can move in a day to calm panic and manipulation. Reports on ESX rules cite a **daily price movement limit of around 10%**, together with trading halts in extreme volatility. **Rules change, so verify the current limit and circuit-breaker thresholds on esx.et before relying on them.**

## Reading a quote on Stocet

Stocet's **Markets** and **Analytics** pages show for each ticker the last price, percent change, a candlestick chart and indicators. Remember:

- The **last price** may be hours old if a share did not trade today.
- A **big percent change on a tiny volume** says little. Look for repeated trades at consistent prices.
- In a young market, a single large order can move a price a lot.

## Key takeaways

1. Price = where the best bid meets the best ask; the spread is a hidden cost.
2. Default sessions: pre-open then continuous trading, 9:00 to 15:00 EAT on weekdays.
3. Daily limits exist; verify the current numbers with ESX.
4. Thin trading makes prices noisy.
""", quiz=[
 ("What is the 'spread'?", ["The broker's fee", "The difference between the best bid and the best ask", "The daily price range", "The dividend yield"], 1),
 ("The best bid is ETB 1,150 and the best ask is ETB 1,180. What does buying at the ask and immediately selling at the bid cost per share?", ["ETB 0", "ETB 30", "ETB 1,150", "ETB 1,180"], 1),
 ("A wide spread usually signals what?", ["A very liquid market", "Few participants and higher trading costs", "A guaranteed profit", "A halted market"], 1),
 ("What is the purpose of a daily price movement limit?", ["To guarantee profits", "To calm extreme moves and limit manipulation", "To set dividends", "To raise fees"], 1),
 ("Why should you be careful with a big percent move on very small volume?", ["It is always a scam", "A single small trade can move a thin share a lot, so it may mean little", "It means the company paid a dividend", "The exchange is closed"], 1),
]),

dict(title="What Does Trading Really Cost?", video="", content="""
## Costs are a return you never see

Every trade has costs, and they come out of **your** return. In a market where you might expect single-digit or low-double-digit percent returns, costs matter.

## Fee components

Based on the ESX fee schedule dated December 2024 and analyses of it (**check the current schedule, since fees can change**):

| Component | Indicative level | Who gets it |
|---|---|---|
| Exchange transaction fee | about **0.36%** of trade value (main market equities) | ESX |
| Regulator fee | about **0.15%** | ECMA |
| Brokerage commission | **capped at 1.6%** for listed equities (your dealer may charge less) | Your dealer |

Together that is on the order of **2% per side** at the cap, though many dealers charge less, and fee schedules are revised.

## The round trip

You pay on the **buy** and again on the **sell**. If costs are about 1% each way, you need the price to rise **over 2%** just to break even, *before* the spread and tax.

*Example:* buy 100 shares at ETB 1,000 (value ETB 100,000). With roughly 0.51% statutory fees plus 1% brokerage, that is about 1.5%, or **ETB 1,510**, to buy and a similar amount to sell: about ETB 3,000 round trip. The share must rise about 3% to break even.

## Other costs

- **The spread** (see the last lesson).
- **Taxes.** Dividend withholding and capital-gains tax rules for Ethiopian investors have been **reformed recently, and sources disagree on current rates and exemptions** (rates of 10% to 15% on dividends and 15% to 30% on capital gains have all been quoted for different years). Do **not** rely on a course for tax rates: check the latest Ministry of Revenues guidance or ask a licensed tax adviser.
- **Account or custody fees**, if your dealer charges them.

## What costs mean for strategy

1. **Frequent trading is expensive.** Each in-and-out costs perhaps 2% to 3%.
2. **Long-term investing** spreads those costs over years.
3. In Stocet's backtester, set **commission and slippage** to realistic values. A strategy that looks great at 0% costs can lose money at 1.5% per side. Try it.

## Key takeaways

1. Costs apply on both buy and sell and compound against you.
2. The statutory part is roughly 0.5% a side; dealer brokerage is capped at 1.6% (verify current rates).
3. Verify tax treatment from official sources; it has changed recently.
4. Trade less, and always include costs in backtests.
""", quiz=[
 ("Why do trading costs matter more for frequent traders?", ["They don't", "Costs are paid on every buy and sell and add up against returns", "Frequent traders pay lower taxes", "Costs only apply to the first trade"], 1),
 ("If total costs are about 1.5% to buy and 1.5% to sell, roughly how much must the price rise just to break even?", ["0%", "1.5%", "3%", "15%"], 2),
 ("Who should you rely on for current tax rates on dividends and capital gains?", ["A friend", "A course or blog from a year ago", "Official guidance or a licensed tax adviser", "The share price"], 2),
 ("In a backtest, what happens if you set commission and slippage to 0%?", ["Results become more realistic", "Results are usually too optimistic", "Nothing changes", "The test cannot run"], 1),
 ("Which fee goes to the exchange itself?", ["Brokerage commission", "Transaction fee", "Income tax", "Account opening fee"], 1),
]),

dict(title="IPOs and How to Read a Prospectus", video="", content="""
## What an IPO is

An **initial public offering (IPO)**, or public offer, is when a company sells shares to the public **for the first time** (the primary market) and then lists them on the exchange so they can be traded. Later sales of new shares by an already-listed company are called **additional public offers** (for example, Gadaa Bank published an additional public offer notification in 2026).

## The prospectus is your most important document

Before an offer, the company publishes a **prospectus** approved by ECMA. It is long, but you do not need to read every page. Focus on:

1. **Who is issuing and why.** What will the money be used for: expansion, technology, repaying debt, or simply letting existing owners sell? Money raised for growth is generally a better sign than owners simply cashing out.
2. **The offer terms.** Price per share, number of shares, minimum and maximum order sizes, the application window, how allocation works if the offer is oversubscribed, and the listing date.
3. **The business.** What does it do, how does it earn money, who are its customers, what is its market position?
4. **Financial statements.** Several years of audited numbers: profit trend, assets, debt, capital. (The Fundamental Analysis course shows how to read them.)
5. **Risk factors.** The section companies must write, and the one most investors skip. Read it.
6. **Governance.** Who is on the board and management? Who are the largest shareholders? Who is the auditor?
7. **Dividend policy.** How much of profit does the company intend to pay out?

## Valuation: is the price fair?

Compare the offer price with the company's **earnings per share (EPS)** and **book value per share**, and with similar companies. A price of 2x book value and 3x earnings means very different things depending on the company's growth and returns. The Fundamentals pages on Stocet calculate P/E and P/B for every listed company so you can compare.

## Common IPO mistakes

- Buying because "everyone is applying," without reading anything.
- Expecting a first-day jump. Some listings rise, some fall; there is no rule.
- Putting too much money into a single offer.

## Where to find prospectuses

Every ESX-listed company's prospectus and later filings are published on its ESX directory page, and Stocet links them from each company's **Filings** tab.

## Key takeaways

1. An IPO sells new shares to the public and lists them for trading.
2. Read at least: use of proceeds, offer terms, financials, risk factors and governance.
3. Judge the price against earnings and book value, not excitement.
""", quiz=[
 ("In an IPO, who receives the money paid by investors for new shares?", ["The ESX", "ECMA", "The issuing company (mainly)", "The securities dealer only"], 2),
 ("Which prospectus section do many investors skip but should read?", ["Cover page", "Risk factors", "Table of contents", "Auditor's address"], 1),
 ("Which 'use of proceeds' is generally a healthier sign?", ["Existing owners simply selling their stake", "Funding expansion or technology", "Unspecified purposes", "Paying the offering's marketing costs only"], 1),
 ("What can you compare an offer price against to judge whether it is fair?", ["The company's logo", "EPS, book value and similar companies' valuations", "The number of branches only", "The founding year"], 1),
 ("Where can you find the ESX-listed companies' prospectuses on Stocet?", ["The Alerts page", "Each company's Filings tab under Fundamentals", "The Education page", "Settings"], 1),
]),

dict(title="Dividends, AGMs and Other Corporate Actions", video="", content="""
## Corporate actions

A **corporate action** is an event decided by a company that affects its shareholders. Stocet tracks them on the **Calendar** page and each company's profile.

### Dividends
A share of profit paid to shareholders, usually per share. Key dates:

- **Declaration date**: the company announces the dividend (usually after the AGM approves it).
- **Record date**: you must be on the shareholder register on this date.
- **Ex-dividend date**: shares bought **on or after** this date do **not** receive the dividend; the price usually drops by about the dividend amount.
- **Payment date**: the money is paid.

**Dividend yield** = annual dividend per share / share price. **Payout ratio** = dividends / net profit. A very high payout can be unsustainable; a very low one may mean the company is reinvesting for growth.

### Annual General Meeting (AGM)
Shareholders meet once a year to approve accounts, elect directors, appoint the auditor and approve dividends. As a shareholder you can vote, directly or by proxy. ESX-listed companies publish an **AGM notice** and, afterwards, the **AGM resolutions** on their ESX page. Recent examples include Bank of Abyssinia's 17th AGM, Gadaa Bank's 5th AGM and Wegagen Bank's 33rd AGM.

### Rights issue
The company offers **existing shareholders** the right to buy new shares, often at a discount, in proportion to their holding. If you ignore it, your ownership percentage is **diluted**.

### Bonus shares (stock dividend)
Free extra shares issued from reserves. Your *percentage ownership* does not change and, in theory, neither does your total value: you own more shares, each worth proportionally less.

### Share capital increases
Banks in Ethiopia regularly increase share capital (partly to meet regulators' capital requirements), which can dilute shareholders who do not participate.

### Closed periods
Before results are published, companies enter a **closed period** (you will see "Closed Period Announcement" among ESX filings), during which insiders may not trade, to prevent unfair advantage.

## What to do as a small investor

1. Check the **Calendar** for upcoming dates before you buy.
2. Read **AGM notices and resolutions**; they are free and show what the company really decided.
3. For rights issues, decide in time; deadlines are short.

## Key takeaways

1. To receive a dividend you must own the share **before** the ex-dividend date.
2. Rights issues and capital increases can dilute you if you do not take part.
3. AGM resolutions and closed-period notices are public: read them.
""", quiz=[
 ("When does a buyer stop being entitled to an upcoming dividend?", ["After the payment date", "On or after the ex-dividend date", "Before the record date", "On the AGM date"], 1),
 ("What is dividend yield?", ["Dividend / share price", "Profit / equity", "Price / earnings", "Dividend / total shares"], 0),
 ("A company issues a rights issue and you do not participate. What happens to your ownership percentage?", ["It increases", "It is diluted (falls)", "It stays the same", "Your shares are cancelled"], 1),
 ("What is a 'closed period'?", ["A time the exchange is closed", "A time before results when insiders may not trade", "The dividend payment period", "A holiday"], 1),
 ("Where can you read what shareholders actually decided at an AGM?", ["Only by attending", "In the AGM resolutions published on the company's ESX page", "On social media only", "It is secret"], 1),
]),

dict(title="Protecting Yourself: Scams and Red Flags", video="", content="""
## A new market attracts fraudsters

When people are excited about a new market and unfamiliar with how it works, scammers move in. Learn the patterns.

## Common scams

1. **"Guaranteed returns."** No legitimate investment guarantees a high return. If someone promises 5% a month, it is a Ponzi or pyramid scheme.
2. **Unlicensed "brokers" and "advisers."** Only ECMA-licensed firms may deal in or advise on securities. Check the licensee list on ECMA's website and Stocet's CMSPs page **before** sending any money.
3. **Fake apps and websites.** Install the ESX Neway app only from its official store listing and your dealer's official channels. Never log in from a link in a message.
4. **Pump-and-dump groups.** Telegram or WhatsApp groups that push a share "to the moon" so insiders can sell to you.
5. **Pre-IPO "allocations."** Anyone selling you a guaranteed allocation of an IPO, or "pre-listing shares," is almost certainly a fraud.
6. **Impersonation.** Messages claiming to be from ECMA, ESX, NBE or your bank asking you to confirm your PIN or send a fee.
7. **Paying personal accounts.** Money for investments must go to the licensed firm's official account, never a personal account or mobile wallet.

## Red-flag checklist

- Pressure to act **now**.
- A return that is **too good** compared with T-bills and bank deposits.
- Vague explanation of how the money is made.
- Asks you to recruit others.
- Not on ECMA's list.
- Refuses to give written terms.

## What to do

- Verify the firm with ECMA (contact details are on ecma.gov.et).
- Never share passwords, PINs or OTP codes.
- Keep records of all transactions.
- If you were scammed, report it to ECMA and the police promptly.
- **Stocet** never asks for your trading password or money, and never executes trades. Be suspicious of anyone claiming otherwise.

## Key takeaways

1. Guaranteed high returns are a scam signal.
2. Verify every firm against ECMA's licensee list first.
3. Only use official apps and pay only official firm accounts.
""", quiz=[
 ("Someone promises a guaranteed 5% return every month. What is it most likely?", ["A safe investment", "A Ponzi or pyramid scheme", "A treasury bill", "An ESX product"], 1),
 ("Where should you verify whether a broker is legitimate?", ["Their WhatsApp status", "ECMA's licensee list", "Their number of followers", "Their business card"], 1),
 ("A group tells you to buy a share 'before it moons'. This is a typical sign of:", ["A research report", "A pump-and-dump", "A dividend announcement", "An AGM"], 1),
 ("Where should money for share purchases be paid?", ["A friend's mobile wallet", "The licensed firm's official account", "A personal bank account", "Cash to an agent"], 1),
 ("Which is true about Stocet?", ["It places trades on your behalf", "It asks for your trading PIN", "It never executes trades or asks for your trading credentials", "It guarantees profits"], 2),
]),
]
