COURSE = dict(
    key="c6",
    title="Macroeconomics & Ethiopia's Financial System",
    duration="1 hour 45 mins",
    description="Understand the big forces behind share prices: inflation, interest rates, the National Bank of Ethiopia's policy, treasury bills, the birr exchange rate and the story of how Ethiopia built a capital market. Updated with the NBE's July 2026 decision.",
)

LESSONS = [
dict(title="Why Macroeconomics Moves Markets", video="", content="""
## The tide and the boats

Individual companies are like boats; the economy is the tide. A good company can still struggle when the tide goes out. Four macro variables matter most for ESX investors.

### 1. Economic growth (GDP)
Faster growth usually means more lending, more deposits, more phone and data use, and higher company profits.

### 2. Inflation
The rate at which prices rise. It affects investors in two opposite ways:

- **It erodes returns.** A 10% return with 13% inflation is a negative real return.
- **It can inflate nominal profits.** Banks and telecoms that can raise prices or lending rates in step with inflation report growing nominal earnings. Watch the *real* growth.

*As of mid-2026:* headline inflation in Ethiopia was **13.4% in May 2026**, up from 11.7% in April, according to reports on the NBE's announcement, and the NBE expected it to ease by December 2026 while remaining in double digits over its forecast window.

### 3. Interest rates
The price of money. When rates rise:

- Banks can charge more on loans, but borrowers struggle more and some loans go bad.
- Safe assets (T-bills, deposits) become more attractive relative to shares, so valuations can fall.
- Companies with debt pay more interest.

### 4. The exchange rate
Ethiopia imports fuel, medicine, machinery and network equipment. A weaker birr raises costs and prices. (See the exchange-rate lesson.)

## The macro dashboard for an ESX investor

| Indicator | Where to look | Why it matters |
|---|---|---|
| NBE policy rate | nbe.gov.et/monetary-policies | Direction of interest rates |
| Inflation | Ethiopian Statistics Service, NBE | Real returns, policy outlook |
| T-bill yields | nbe.gov.et/treasury-bills | Risk-free alternative to shares |
| Exchange rate | nbe.gov.et/exchange | Costs, imports, dollar profits |
| Credit growth | NBE statistics | Fuel for bank profits |

## Reading the signs

1. **Rising rates + high inflation** usually pressure valuations of shares but can boost bank margins in the short term.
2. **Falling inflation + stable currency** is generally supportive for share prices over time.
3. **Policy changes** (like removing credit caps) can quickly change what banks can earn, so read the announcements.

## Key takeaways

1. Macro variables (growth, inflation, interest rates, FX) shape every company's results.
2. Always think in **real** (after inflation) terms.
3. Follow NBE announcements; they move the whole banking sector.
""", quiz=[
 ("A 10% investment return with 13% inflation is a:", ["Positive real return", "Negative real return", "Zero return", "Risk-free return"], 1),
 ("Which is generally true when interest rates rise?", ["Safe assets become less attractive", "Safe assets such as T-bills become more attractive compared with shares", "All shares rise", "Loans become cheaper"], 1),
 ("Why does the exchange rate matter for companies like Ethio Telecom?", ["It sets dividends", "They import equipment priced in dollars", "It changes the ticker", "It has no effect"], 1),
 ("Where can you find the NBE's policy rate?", ["The ESX trading screen", "nbe.gov.et monetary policy pages", "A company prospectus", "Stocet Settings"], 1),
 ("Why look at 'real' growth in profits?", ["Nominal profits can rise just because of inflation", "Real is always bigger", "Nominal is illegal", "Taxes depend on it"], 0),
]),

dict(title="The National Bank of Ethiopia and Monetary Policy", video="", content="""
## What the NBE does

The **National Bank of Ethiopia (NBE)** is the central bank. It:

- Sets **monetary policy**, aimed at stable prices.
- **Supervises commercial banks**: six of the seven ESX-listed companies.
- Runs **treasury-bill auctions** for the government.
- Manages **foreign-exchange** policy and reserves.

## The policy rate and the corridor

Since 2024 the NBE has used an interest-rate-based framework with a **policy rate** (the National Bank Rate, NBR) and an **interest-rate corridor**. The commercial banks' overnight borrowing and lending rates sit around the policy rate:

- **Policy rate**: the benchmark.
- **Corridor**: +/- 3 percentage points. Banks can deposit surplus funds with the NBE at the lower bound and borrow overnight at the upper bound.

*Latest decision (13 July 2026):* the NBE's Monetary Policy Committee **raised the policy rate by one percentage point to 16%**, its first change since the rate was introduced in 2024, citing renewed inflation pressure from the global oil price shock. With the unchanged corridor, the overnight deposit rate is 13% and the lending rate 19%. The next meeting was scheduled for the end of September 2026: **always check the NBE website for the latest decision**, because this may have changed.

## Credit controls and reserves

In the same decision the NBE:

- **Removed the cap on commercial-bank credit growth** (previously 24%).
- Introduced a **targeted reserve-requirement tool** that lets it impose bank-specific reserve ratios on banks lending aggressively relative to their deposits.

*What this means for bank investors:* removing the credit cap lets banks lend more freely (potentially more profit), but a reserve tool means regulators can restrict banks that grow too fast. Bank growth in FY2025/26 (assets up 30% to 65%) happened under the credit cap, so the picture may change.

## Treasury bills: the risk-free benchmark

The NBE auctions T-bills of different maturities (for example 28, 91, 182 and 364 days). The **yield** is the return you earn. Why investors care:

- T-bill yields are the **risk-free rate**: the minimum return you should demand from a risky investment.
- If T-bills pay more than a share's earnings yield (1 / P/E) or dividend yield, shares must compensate through expected growth.
- Large government borrowing through T-bills can crowd out private lending.

Check the latest yields on nbe.gov.et/treasury-bills.

## Reading an NBE announcement

1. Did the policy rate change, and by how much (compared to expectations)?
2. What does it say about **inflation** forecasts?
3. Any changes to **reserve requirements, credit limits or FX rules**?
4. When is the **next meeting**?

## Key takeaways

1. The NBE sets the policy rate and supervises the banks that dominate the ESX.
2. In July 2026 it raised the policy rate to 16% and removed the credit-growth cap.
3. T-bill yields are the risk-free benchmark for every investment decision.
""", quiz=[
 ("In July 2026 the NBE moved its policy rate to:", ["14%", "15%", "16%", "20%"], 2),
 ("What does a +/- 3 point corridor around a 16% policy rate imply for the overnight lending rate?", ["13%", "16%", "19%", "22%"], 2),
 ("What is the 'risk-free rate' usually based on?", ["Bank deposit rates only", "Treasury bill yields", "Share dividends", "Inflation"], 1),
 ("What did the NBE remove in its July 2026 decision?", ["The policy rate", "The cap on commercial-bank credit growth", "T-bill auctions", "The ESX"], 1),
 ("Where should you check for the NBE's latest decision?", ["A friend's message", "nbe.gov.et", "An old article", "A prospectus"], 1),
]),

dict(title="The Birr and the Exchange Rate", video="", content="""
## From a managed rate to a market rate

For years, Ethiopia's exchange rate was **tightly managed**, and a large gap developed between the official rate and the parallel (black) market rate. In **July 2024** the NBE moved to a **market-based foreign-exchange system**, letting the birr's value be determined largely by supply and demand in the interbank market. The birr fell sharply at the start (reported at roughly 30% against the dollar as the system changed) and has continued to adjust since.

## Why it matters to ESX investors

**Winners and losers:**

- **Importers and companies with dollar costs** (like telecoms buying network equipment, or manufacturers buying raw materials) see costs rise in birr terms when the birr weakens.
- **Exporters and companies with dollar income** gain in birr terms.
- **Banks** can earn from foreign-exchange trading and trade finance. Several banks' reports show large swings in "net foreign exchange income" after the float. (Bank of Abyssinia, for example, moved from a small FX loss in FY2024/25 to ETB 1.68 bn of FX income in FY2025/26.)
- **Dollar-denominated debt** becomes more expensive in birr.
- **Foreign investors** worry about converting returns back to dollars.

## Reading FX in your analysis

1. **Where does the company earn and spend?** Birr revenue and dollar costs mean pressure when the birr falls.
2. **FX income can be volatile.** Do not assume a one-time boost in FX gains will repeat.
3. **Inflation link:** weaker birr feeds inflation through imported goods, which is why the NBE watches both together.

## Recent related policy changes

In July 2026 the NBE also reduced the **foreign-exchange commission** (from 2.5% to 1.5%) and the **export-proceeds surrender requirement** (from 50% to 30%), per reports on that decision. Such changes affect how much foreign currency banks and exporters get to keep and trade.

Check the latest official rates on nbe.gov.et/exchange; they are published every working day.

## Key takeaways

1. Since July 2024 the birr is market-determined, so it moves.
2. A weaker birr hurts firms with dollar costs and helps those with dollar income.
3. FX gains in a bank's accounts can be one-off; look for what repeats.
""", quiz=[
 ("In July 2024 the NBE moved Ethiopia to:", ["A fixed peg", "A market-based exchange-rate system", "The US dollar", "A currency board"], 1),
 ("Who typically loses when the birr weakens?", ["Exporters with dollar income", "Importers with dollar costs", "Holders of dollars", "Nobody"], 1),
 ("Why be careful when a bank reports a large foreign-exchange gain?", ["FX gains are taxed twice", "They may be one-off and not repeat", "They are illegal", "They reduce capital"], 1),
 ("Where are official exchange rates published?", ["nbe.gov.et/exchange", "A prospectus", "ESX order book", "ECMA licence list"], 0),
 ("A weaker birr tends to feed inflation mainly through:", ["Higher dividends", "More expensive imported goods", "Lower interest rates", "More T-bills"], 1),
]),

dict(title="How Ethiopia Built a Capital Market", video="", content="""
## The timeline

| When | What happened |
|---|---|
| 2021 | The Capital Market Proclamation (No. 1248/2021) creates the legal framework and the Ethiopian Capital Market Authority (ECMA); a new Commercial Code (No. 1243/2021) is also enacted |
| 2024 | The NBE moves to a market-based foreign exchange system (July); ECMA begins licensing, starting with securities investment advisors |
| 2025 (10 January) | **Wegagen Bank** becomes the first company listed on the ESX |
| 2025 (March) | ECMA licenses the first investment banks and securities dealer: CBE Capital, Wegagen Capital and Ethio-Fidelity Securities |
| 2025 (13 June) | **Gadaa Bank** lists |
| 2026 | **Awash Bank** (23 April), **Ethio Telecom** (26 May), **Abay Bank** (25 June), **Bank of Abyssinia** (28 July) list; **Sidama Bank** follows on 28 September |
| October 2026 | ECMA's directory shows 20 licensed service providers (7 investment banks, 2 securities dealers, 11 securities investment advisors) |

## Why it was built

- **Mobilise savings** that were not reaching productive investment.
- **Give companies an alternative to bank debt**, and a way to raise long-term capital.
- **Let citizens share in growth** by owning shares of national companies.
- **Develop the banks**: share listings give banks a stronger capital base and public accountability.

## The architecture

- **ECMA** regulates and licenses.
- **ESX** runs the market (equities and, with fixed-income trading rules, debt).
- **Licensed intermediaries** (investment banks, securities dealers, advisors) serve investors and issuers.
- **ESX Neway** is the trading app used with licensed intermediaries.
- **Disclosure**: listed companies file prospectuses, financial statements, closed-period notices, AGM notices and material announcements on the ESX site, which Stocet collects in each company's Filings tab.

## Where it is heading

The market is at an early stage. Plausible developments, though none are guaranteed:

1. **More listings**: other banks, insurers, and companies across other sectors have registered with ECMA but not yet listed.
2. **Debt market development**: government and corporate bonds, building on the T-bill market.
3. **Funds and ETFs**: collective investment products for small investors.
4. **Greater liquidity**: more investors, market makers and trading.
5. **Better data**: more disclosures, more research (a gap Stocet exists to help fill).

## What it means for you

- You are an **early participant**: potential rewards, and the risks of a new market (thin trading, limited track record, concentration in banks, rapid rule changes).
- **Learn the rules** and verify them on official sites; they are evolving.
- Be patient: markets mature over years, not months.

## Key takeaways

1. Ethiopia's capital market has moved from law (2021) to a working exchange (Jan 2025) with seven listed companies by late 2026.
2. ECMA regulates, ESX operates, licensed intermediaries serve you.
3. Early-stage markets offer opportunity and extra risk: verify rules and diversify.
""", quiz=[
 ("Which company was the first listed on the ESX?", ["Awash Bank", "Wegagen Bank", "Ethio Telecom", "Sidama Bank"], 1),
 ("How many ECMA-licensed capital market service providers were listed in October 2026?", ["8", "12", "20", "50"], 2),
 ("Which body regulates and licenses market participants?", ["ESX", "ECMA", "NBE", "Ministry of Trade"], 1),
 ("Which of these is a risk specific to an early-stage market?", ["Thin trading and rapid rule changes", "Guaranteed liquidity", "No information at all", "No intermediaries"], 0),
 ("Where does Stocet collect companies' prospectuses and announcements from?", ["Social media", "The ESX website's company directory pages", "Anonymous tips", "Bank branches"], 1),
]),
]
