COURSE = dict(
    key="c3",
    title="Fundamental Analysis for ESX Stocks",
    duration="3 hours",
    description="Learn to judge a company the way professional analysts do: read the income statement, balance sheet and cash flow; understand how banks and a telecom actually make money; use EPS, P/E, P/B, ROE and ROA; compare the six listed banks with real FY2025/26 audited numbers; and spot red flags in an audit report.",
)

LESSONS = [
dict(title="Reading the Three Financial Statements", video="", content="""
## Why fundamentals matter

A share is a claim on a business. **Fundamental analysis** asks: *what is this business worth, and is the price fair?* The raw material is the company's published financial statements, which ESX-listed companies must file. On Stocet you can open any company under **Fundamentals** and see the figures and the original PDFs.

## 1. The income statement (profit or loss)

It shows performance **over a period** (a financial year). Read it top to bottom:

**Revenue (or total operating income) - expenses = profit**

- **Operating income / revenue**: what the business earned from its activities.
- **Operating expenses**: staff, depreciation, rent, and other costs of running the business.
- **Impairment charges**: money set aside for losses (very important for banks).
- **Profit before tax (PBT)**, then **tax**, then **Profit after tax (PAT) / net profit**: what belongs to shareholders.
- **Earnings per share (EPS)**: net profit divided by the number of shares.

## 2. The balance sheet (statement of financial position)

A **snapshot at one date** (for most ESX companies, 30 June). It always balances:

**Assets = Liabilities + Equity**

- **Assets**: what the company owns or is owed (cash, loans it has made, buildings, equipment).
- **Liabilities**: what it owes (customer deposits for a bank, borrowings, payables).
- **Equity**: what is left for the owners (share capital, reserves, retained profit). Also called **net assets** or **book value**.

## 3. The cash flow statement

Profit is an opinion; cash is a fact. This statement shows actual cash moving through **operating**, **investing** and **financing** activities. A healthy company generates cash from operations that is broadly in line with its profit over time.

## How they connect

Net profit from the income statement (after any dividends) flows into **retained earnings** in equity on the balance sheet, and the cash side of everything appears in the cash flow statement.

## Where the numbers hide

- **Units:** Ethiopian statements usually show **Birr '000** (thousands). 29,341,603 in Birr '000 means ETB **29.3 billion**.
- **Consolidated vs. separate:** a group may publish statements for the parent alone ("separate") and for the group including subsidiaries ("consolidated"). Compare like with like.
- **Comparatives:** every statement shows the current year next to the previous year. Always compare them.

## Key takeaways

1. Income statement = performance over a period; balance sheet = position at a date; cash flow = actual cash.
2. Assets = Liabilities + Equity.
3. Check units, consolidated vs. separate, and the prior-year column.
""", quiz=[
 ("Which statement is a snapshot at a single date?", ["Income statement", "Balance sheet", "Cash flow statement", "AGM minutes"], 1),
 ("What is the accounting identity?", ["Assets = Equity - Liabilities", "Assets = Liabilities + Equity", "Profit = Assets + Equity", "Equity = Revenue - Cost"], 1),
 ("A statement shows profit after tax of 29,341,603 in Birr '000. What is that in ETB?", ["ETB 29,341,603", "ETB 29.3 million", "ETB 29.3 billion", "ETB 2.9 billion"], 2),
 ("Why look at the cash flow statement as well as profit?", ["Profit can include non-cash items; cash is a fact", "It is shorter", "It shows the dividend", "It replaces the balance sheet"], 0),
 ("What does 'separate' mean on a bank's statements?", ["Statements for the parent company alone", "Statements prepared separately by the regulator", "Statements for only one branch", "Unaudited statements"], 0),
]),

dict(title="How a Bank Makes Money: Reading Bank Accounts", video="", content="""
## Banks are different

Six of the seven ESX-listed companies are banks, and bank statements look unlike a normal company's. A bank's "raw material" is **money itself**: it borrows from depositors and lends to borrowers.

## The bank income statement

Using **Awash Bank's audited FY2025/26 separate statements** as a real example (ETB):

| Line | Amount | What it means |
|---|---|---|
| Interest income | 54.17 bn | Earned on loans and securities |
| Interest expense | (15.34 bn) | Paid to depositors and lenders |
| **Net interest income (NII)** | **38.83 bn** | The core spread: the bank's main profit engine |
| Fees and commission | 22.91 bn | Charges for payments, trade finance, guarantees, etc. |
| Other operating income | 7.43 bn | FX gains and other |
| **Total operating income** | **69.17 bn** | |
| Impairment on financial assets | (1.25 bn) | Provisions for bad loans |
| Personnel and other costs | (about 29.6 bn) | Staff, branches, technology |
| **Profit before tax** | **38.31 bn** | |
| Tax | (8.97 bn) | |
| **Profit after tax** | **29.34 bn** | |

**Net interest income** is the biggest driver, and **fee income** shows how diversified a bank is. Compare banks by *how much of their income comes from fees* and how fast each grows.

## The bank balance sheet

Awash, 30 June 2026 (separate):

- **Assets** 624.1 bn: loans 268.5 bn (the largest earning asset), cash and bank balances 234.7 bn, government and other securities ~68 bn.
- **Liabilities** 533.7 bn: customer deposits 469.5 bn (88% of liabilities).
- **Equity** 90.4 bn.

The bank's assets are mostly **other people's money** (deposits). Equity is the thin cushion that absorbs losses, which is why regulators set **minimum capital** rules and why you should look at the equity-to-assets ratio (Awash: 14.5%).

## Bank-specific metrics

- **Loan-to-deposit ratio (LDR)** = loans / deposits. Awash: 268.5 / 469.5 = **57%**. Higher means more of deposits are lent out (more earning power, less liquidity buffer).
- **Net interest margin** = NII / earning assets.
- **Cost-to-income ratio** = operating costs / operating income. Lower is more efficient.
- **Non-performing loan (NPL) ratio**: the share of loans not being repaid. Rising NPLs are the classic sign of trouble. *Always check the notes for this.*
- **Capital adequacy ratio (CAR)**: regulatory capital relative to risk-weighted assets. Must stay above the NBE minimum.
- **Customer deposits growth** shows whether customers trust the bank.

## Interest-free (Islamic) banking

Many Ethiopian banks have **interest-free windows**. Bank of Abyssinia, for example, reports interest-free customer deposits (about ETB 40.8 bn of its ETB 318.8 bn total deposits in FY2025/26). In Stocet, "customer deposits" includes both.

## Key takeaways

1. A bank earns the **spread** between what it charges borrowers and pays depositors (NII), plus fees.
2. Equity is a thin cushion on a balance sheet full of deposits; capital strength matters.
3. Always check loan quality (NPLs and impairments), not just profit.
""", quiz=[
 ("What is net interest income?", ["Total fees earned", "Interest income minus interest expense", "Profit after tax", "Deposits minus loans"], 1),
 ("Awash Bank's FY2025/26 loans were ETB 268.5 bn and deposits ETB 469.5 bn. What is the loan-to-deposit ratio (about)?", ["25%", "57%", "88%", "175%"], 1),
 ("What is the main bank balance-sheet liability?", ["Share capital", "Customer deposits", "Retained earnings", "Fixed assets"], 1),
 ("Which sign most clearly warns of trouble in a bank's loan book?", ["Rising customer deposits", "Rising non-performing loans", "More branches", "Higher fee income"], 1),
 ("Why does a bank's equity-to-assets ratio matter?", ["It sets the share price", "Equity is the cushion that absorbs losses on deposit-funded assets", "It determines the dividend date", "It shows the number of customers"], 1),
]),

dict(title="The Key Ratios: EPS, P/E, P/B, ROE, ROA", video="", content="""
## Why ratios

Absolute numbers do not compare across companies of different sizes. **Ratios** put them on a common scale. Stocet's **Fundamentals** and **Screener** pages calculate these for every listed company.

## Per-share and valuation ratios

**EPS (Earnings per share)** = net profit / number of shares (weighted over the year).

> **Watch out in Ethiopia:** many companies state EPS as a **percentage of par value** or per Birr 1 of capital. Bank of Abyssinia reports FY2025/26 EPS as "77%", which means **ETB 770 per share** at a par value of ETB 1,000. Gadaa Bank reports "Birr 0.52", meaning ETB 520 per ETB 1,000 share. **Stocet normalises every EPS to ETB per share.**

**P/E (Price-to-earnings)** = price / EPS. How many years of current earnings you pay for. A low P/E can mean cheap *or* that the market expects earnings to fall.

**P/B (Price-to-book)** = price / book value per share, where book value per share = equity / shares. P/B below 1 means you can buy the net assets for less than their accounting value.

**Dividend yield** = dividend per share / price.

*Worked example (illustrative price):* Awash Bank FY2025/26: EPS = ETB 858; equity = ETB 90.39 bn; ESX-listed shares = 37,896,928.
- Book value per share = 90.39 bn / 37.897 m = **ETB 2,385**.
- If the price were ETB 2,000: **P/E = 2,000 / 858 = 2.3x**; **P/B = 2,000 / 2,385 = 0.84x**.

(Check the live price on Stocet; the numbers above show the method, not a recommendation.)

## Profitability ratios

- **Net margin** = net profit / total operating income. Awash: 29.34 / 69.17 = 42%.
- **ROE (Return on equity)** = net profit / equity. How much profit each birr of owners' money generates. Awash: 29.34 / 90.39 = **32.5%** (on year-end equity).
- **ROA (Return on assets)** = net profit / total assets. Awash: 29.34 / 624.1 = **4.7%**. For banks, ROA above 2% is considered strong.

## Safety ratios

- **Debt-to-equity** = total liabilities / equity. For banks this is naturally high (Awash: 533.7 / 90.4 = 5.9x); for non-banks, it is a key leverage test.
- **Equity-to-assets** = equity / assets. Bigger is safer.
- **Loan-to-deposit ratio** (banks).

## How to use them

1. **Compare against peers** (banks with banks) and against the same company in prior years.
2. **No ratio works alone.** A tiny P/E on collapsing earnings is a trap; a great ROE funded by extreme leverage is risky.
3. **Use consistent and current numbers.** A P/E that uses last year's EPS with this year's price is stale if earnings have grown.

## Key takeaways

1. Convert EPS to ETB per share before using it; Stocet does this for you.
2. P/E tells you what you pay for earnings; P/B what you pay for net assets.
3. ROE and ROA measure how efficiently a company turns capital into profit.
""", quiz=[
 ("Bank of Abyssinia reports EPS of '77%' on a par value of ETB 1,000. What is EPS per share?", ["ETB 0.77", "ETB 77", "ETB 770", "ETB 7,700"], 2),
 ("A share costs ETB 2,000 and EPS is ETB 858. What is the P/E ratio (about)?", ["0.43x", "2.3x", "23x", "42.9x"], 1),
 ("A P/B ratio below 1 means:", ["The company is making losses", "The price is below book value per share", "The dividend is zero", "The share is overvalued"], 1),
 ("ROE is calculated as:", ["Net profit / total assets", "Net profit / equity", "Equity / assets", "Price / earnings"], 1),
 ("Why is a very low P/E not automatically a bargain?", ["It always means a dividend is coming", "The market may expect earnings to fall", "It means the price is high", "It proves the company is safe"], 1),
]),

dict(title="Comparing the ESX Banks (FY2025/26)", video="", content="""
## A like-for-like look

Below are key figures from the banks' **audited financial statements for the year ended 30 June 2026**, with ratios calculated by Stocet. Amounts in ETB billions unless noted. *Awash and Bank of Abyssinia figures are from separate (bank-only) statements; Wegagen and Gadaa are consolidated; Abay is separate.* Sidama's FY2025/26 audited statements had not yet been published when this course was written.

| Bank | Total assets | Net profit | Profit growth | ROE | ROA | Loans / deposits | Equity / assets | EPS (ETB) |
|---|---|---|---|---|---|---|---|---|
| Awash (AWAB) | 624.1 | 29.34 | +56.8% | 32.5% | 4.7% | 57% | 14.5% | 858 |
| Bank of Abyssinia (BOAX) | 381.5 | 11.54 | +58.1% | 29.2% | 3.0% | 78% | 10.4% | 770 |
| Abay (ABAYB) | 119.7 | 4.55 | +51.5% | 25.9% | 3.8% | 66% | 14.7% | 525 |
| Wegagen (WGBX) | 111.3 | 3.15 | +13.6% | 18.2% | 2.8% | 88% | 15.6% | 344 |
| Gadaa (GDAB) | 16.8 | 0.76 | +121.5% | 27.4% | 4.5% | 50% | 16.5% | 520 |

*Sidama (SIDAX), FY2024/25: assets ETB 5.52 bn, net profit ETB 72.3 million, EPS ETB 75.*

## What the table tells you

1. **Scale:** Awash is more than 5x the size of Wegagen by assets, and 37x Gadaa. Scale gives pricing power and funding advantages, but small banks can grow faster.
2. **Growth:** Most banks grew profit by 50% or more. Wegagen's profit grew 13.6% even though its assets grew 32%: *growth in assets did not translate into profit growth.* That is worth a question.
3. **Profitability:** Awash leads in ROE (32.5%) and ROA (4.7%). Gadaa's ROA (4.5%) is high for a small young bank.
4. **Lending intensity:** Wegagen lends 88 cents of each deposit birr; Gadaa only 50 cents. A high ratio boosts income but leaves less liquidity if deposits are withdrawn.
5. **Capital:** Equity-to-assets between 10% and 17% looks comfortable, but Bank of Abyssinia is the thinnest (10.4%); it also grew assets 33%.

## What is *not* in the table

- **Loan quality (NPLs)** and provisioning policy.
- **Dividend** history and policy.
- **Price:** A great bank can be an expensive share. Compare ROE with P/B: a bank earning 30% on equity deserves a higher P/B than one earning 15%.
- **One year only.** Banks in Ethiopia have benefited from high inflation, high lending rates and strong deposit growth. Margins can change when policy changes (for example, the NBE raised its policy rate to 16% in July 2026 and removed credit growth caps).

## A simple framework

For each bank, rank on: **growth, profitability (ROE), safety (equity/assets, NPL), efficiency (cost-to-income) and price (P/E, P/B)**. The best investment is usually a *good-enough business at a sensible price*, not the single best business at any price.

Use Stocet's **Fundamentals** and **Screener** to reproduce and extend this table with live prices.

## Key takeaways

1. Compare banks on growth, ROE, ROA, lending intensity and capital together.
2. Profit growth lower than asset growth is a flag to investigate.
3. A strong bank is not automatically a cheap share: always add price.
""", quiz=[
 ("Which bank had the highest ROE in FY2025/26 among those in the table?", ["Wegagen", "Abay", "Awash", "Gadaa"], 2),
 ("Wegagen grew assets 32% but profit only 13.6%. What should an analyst do?", ["Ignore it", "Investigate why growth did not turn into profit", "Assume it is a bargain", "Sell immediately"], 1),
 ("Which bank had the highest loan-to-deposit ratio in the table?", ["Gadaa", "Awash", "Wegagen", "Abay"], 2),
 ("Why must price be added to a comparison of banks' fundamentals?", ["Price is irrelevant", "A strong bank can still be an expensive share", "Price decides profit", "Banks have no price"], 1),
 ("Which important area is NOT shown in the comparison table?", ["Total assets", "Return on equity", "Loan quality (NPLs)", "EPS"], 2),
]),

dict(title="Analysing Ethio Telecom: A Non-Bank Example", video="", content="""
## A very different business

Ethio Telecom (TELE) sells connectivity, infrastructure sharing and mobile money (telebirr). It is **capital-intensive**: it must keep spending on networks. That changes what you look for.

## Audited history (Birr, from the prospectus's audited statements)

| Year to 30 June | Revenue | EBITDA | Operating profit | Profit after tax | Total assets | Equity |
|---|---|---|---|---|---|---|
| 2023 (restated) | 71.5 bn | 37.6 bn | 31.1 bn | 18.0 bn | 199.7 bn | 104.3 bn |
| 2024 | 91.4 bn | 41.2 bn | 32.9 bn | 19.0 bn | 214.2 bn | 113.5 bn |

And **FY2025/26** revenue was reported at about **ETB 215.8 billion**, up 33.2% year on year, with about 90 million customers and 60.6 million telebirr customers. Audited FY2025/26 profit and balance-sheet numbers were not yet on the ESX page when this lesson was written, so Stocet shows revenue only for that year.

## Ratios that fit a telecom

- **EBITDA margin** = EBITDA / revenue. FY2024: 41.2 / 91.4 = **45%**. EBITDA is earnings before interest, tax, depreciation and amortisation. It is widely used for telecoms because depreciation of networks is large.
- **Operating margin** = operating profit / revenue = 36% in FY2024.
- **Net margin** = PAT / revenue = 20.8% in FY2024.
- **Revenue growth**: 27.7% in FY2024 (71.5 to 91.4), then 33.2% reported for FY2025/26.
- **Asset intensity**: property, plant and equipment of ETB 141.5 bn were about two-thirds of total assets.
- **Equity-to-assets** = 53%; borrowings only about 17.7 bn: a conservative balance sheet.
- **ROE** = 19.0 / 113.5 = about 17%.

## Valuation

Ethio Telecom's audited capital is ETB 100 bn at a par value of ETB 100, which implies **about 1 billion shares**. EPS for FY2024 = 19.0 bn / 1 bn = **ETB 19.01**. If the share price were, say, ETB 430, the **P/E on that old EPS would be about 22.6x**. But revenue and profit have grown since, so the P/E on *current* EPS may be lower. This is the **stale-EPS problem**: always check which year's earnings a P/E uses, and update when new audited accounts arrive.

> **Note on ESX listing figures.** ESX lists only the shares offered to the public (10,690,149 shares sold at ETB 300 each, ETB 3.2 bn). The company's total number of shares is far larger, and market capitalisation should use **all** shares, not just the listed tranche. Stocet uses the total.

## What else to check

- **Subscriber growth and ARPU** (average revenue per user).
- **Capital spending** needs and how much cash is left after investing.
- **Dividend policy.**
- **Foreign-exchange exposure**: network equipment is bought in dollars, revenue is in birr.
- **Regulation and competition**: new market entrants.

## Key takeaways

1. Different industries need different ratios: telecoms use EBITDA and capital-intensity measures.
2. Use the *total* number of shares for market cap and EPS, not just those listed in the IPO.
3. Beware stale earnings in a P/E; the market looks forward.
""", quiz=[
 ("What is EBITDA?", ["Earnings before interest, taxes, depreciation and amortisation", "Equity before income, tax, debt and assets", "Expenses before income tax", "Earnings by item, tax and dividends"], 0),
 ("Ethio Telecom's FY2024 EBITDA was ETB 41.2 bn on revenue of ETB 91.4 bn. The EBITDA margin is about:", ["12%", "25%", "45%", "70%"], 2),
 ("Why should market cap use total shares rather than only the shares listed in the IPO?", ["It is a ESX rule", "Market cap is the value of the whole company, not just the offered portion", "Listed shares are larger", "Fewer shares mean higher value"], 1),
 ("What is the 'stale EPS' problem?", ["EPS is always wrong", "A P/E computed with old earnings may overstate valuation if profit has grown", "EPS cannot be calculated for telecoms", "EPS changes daily"], 1),
 ("Why is a telecom called 'capital-intensive'?", ["It has many employees", "It must keep investing heavily in networks and equipment", "It pays high dividends", "It has no debt"], 1),
]),

dict(title="Red Flags and Reading the Audit Report", video="", content="""
## Numbers can mislead

Accounts follow rules (IFRS), but management makes many judgements. Good analysts look for signs that the picture may be rosier than reality.

## The auditor's report

Every audited statement includes an **independent auditor's report**. Read the first page:

- **Unqualified ("clean") opinion**: the statements give a true and fair view. This is what you want.
- **Qualified opinion**: the auditor disagrees with or cannot verify something material. A warning.
- **Adverse opinion / disclaimer**: serious; avoid.
- **Emphasis of matter / material uncertainty about going concern**: the auditor flags doubt that the company can continue.
- **Key audit matters (KAMs)**: the areas that took the most auditor judgement. For banks this is nearly always **impairment of loans (expected credit losses)**. Read what the auditor tested.

## Red flags

1. **Profit growing, operating cash flow not.** Profits may not be turning into cash.
2. **Rapid loan growth with thin provisions.** Growth that looks great today can turn into bad loans later.
3. **Frequent restatements.** Prior-year figures changed ("as restated") more than occasionally.
4. **Unexplained related-party transactions** (loans to directors, shareholders or their companies).
5. **Rising non-performing loans, or large write-offs.**
6. **Late filing**, or changes of auditor without a clear reason.
7. **Dividends that exceed cash generation.**
8. **Share-capital increases at high frequency**, diluting existing owners.
9. **Complex structures** you cannot understand.
10. **Statements that do not foot**: the balance sheet should balance. (Small typographical errors do occur in scanned PDFs; a large unexplained difference deserves a question to the company.)

## Practical checks on ESX filings

- Open the **audited financial statements** from the company's Filings tab on Stocet. They are scanned PDFs, so do not expect to copy numbers from them.
- Compare **current vs. prior year** for each line, and read the **notes** on loans, impairment and related parties.
- Read the **Directors' report** for what management says about the year and the dividend proposal.
- Compare the **provisional announcement** (made soon after year-end) with the final audited numbers. Differences tell you about management's forecasting and the audit's impact. For example, Awash Bank's provisional FY2025/26 statement reported total assets of ETB 622.7 bn; the audited separate statements show ETB 624.1 bn.

## A reminder on limits

A clean audit does not mean a good investment; it means the numbers are fairly presented. Fundamental analysis reduces mistakes but cannot remove risk, especially in a young market with limited disclosure history.

## Key takeaways

1. Read the auditor's opinion first and the key audit matters next.
2. Check cash flow against profit, loan growth against provisions, and prior-year restatements.
3. Compare provisional announcements with audited accounts.
""", quiz=[
 ("What does an 'unqualified' audit opinion mean?", ["The auditor is not qualified", "The statements give a true and fair view", "The company paid no tax", "The company has no debt"], 1),
 ("For banks, the most common key audit matter is:", ["Office rent", "Impairment (expected credit losses) on loans", "Staff training", "Advertising"], 1),
 ("Which is a classic red flag?", ["Profit and operating cash flow both rising", "Rapid loan growth with thin provisions", "Clean audit opinion", "Stable auditor"], 1),
 ("Why compare a provisional announcement with the audited accounts?", ["They are always identical", "Differences show forecasting quality and audit adjustments", "Provisional numbers are illegal", "To find the dividend"], 1),
 ("What does a clean audit tell you?", ["The company is a good investment", "The numbers are fairly presented, not that the investment is good", "The share is cheap", "There is no risk"], 1),
]),
]
