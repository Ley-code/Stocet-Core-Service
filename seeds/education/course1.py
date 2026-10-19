COURSE = dict(
    key="c1",
    existing_id="9f636876-646e-44c2-80cb-362501c05bd4",
    title="Capital Markets 101: From Savings to Shares",
    duration="2 hours",
    description="Start here. Learn how capital markets turn savings into investment, what stocks, bonds and T-bills are, who the key players in Ethiopia's new market are, and the two ideas behind every investing decision: risk and the time value of money. No finance background needed.",
)

LESSONS = [
dict(existing_title="What is a Capital Market?", title="What is a Capital Market?", video="https://youtu.be/SZZG9UIRaXs", content="""
## Money that goes to work

Every economy has people with **spare money** (households, companies with profits, pension funds) and people who **need money** to build things (a bank expanding its branches, a telecom building towers, a government building roads). A **capital market** is the organised marketplace that connects the two groups.

Instead of money sitting under a mattress, it is lent or invested in return for a promise:

- **Lend** your money and you receive **interest** and your money back later. This is **debt** (for example a bond or a treasury bill).
- **Invest** your money and you become a **part-owner** of the business. You share in its profits and growth, and its risks. This is **equity** (a share or stock).

## Primary and secondary markets

| | Primary market | Secondary market |
|---|---|---|
| What happens | A company or government **issues new securities** for the first time | Investors **trade existing securities** with each other |
| Who gets the money | The issuer | The seller (not the issuer) |
| Ethiopian example | An IPO or public offer, such as the public offering of Ethio Telecom shares | Buying and selling Awash Bank (AWAB) shares on the ESX |

The secondary market matters even though the company receives nothing from it, because it gives investors **liquidity**: the confidence that they can sell later. People are far more willing to invest in the primary market when they know a secondary market exists.

## Capital market vs. money market

- **Money market**: short-term borrowing and lending (up to one year), such as treasury bills.
- **Capital market**: longer-term financing, such as shares and bonds.

## Why it matters in Ethiopia

For decades, savings in Ethiopia mostly sat in bank deposits, and companies mostly funded growth from bank loans or retained profits. With the **Ethiopian Capital Market Authority (ECMA)** and the **Ethiopian Securities Exchange (ESX)**, ordinary people can now own shares of Ethiopian companies, and companies can raise money by selling shares instead of only borrowing.

## Key takeaways

1. A capital market connects people who have money with people who need it.
2. Debt means you are a **lender**; equity means you are an **owner**.
3. New securities are sold in the **primary** market; the **secondary** market is where investors trade afterwards.
4. Liquidity (the ability to sell) is what makes people willing to invest.
""", quiz=[
 ("When an investor buys Awash Bank shares from another investor on the ESX, which market is this?", ["Primary market", "Secondary market", "Money market", "Foreign exchange market"], 1),
 ("What does holding equity (a share) make you?", ["A lender to the company", "A part-owner of the company", "A guarantor of its debt", "An employee"], 1),
 ("In the secondary market, who receives the money when you buy a share?", ["The company that issued it", "The Ethiopian Capital Market Authority", "The investor who sold it to you", "The National Bank of Ethiopia"], 2),
 ("Why does a secondary market make people more willing to invest in IPOs?", ["It guarantees profits", "It lets them sell later (liquidity)", "It removes all risk", "It sets the dividend"], 1),
 ("Which of these is a money-market instrument (under one year)?", ["A 10-year corporate bond", "A share of Ethio Telecom", "A 91-day treasury bill", "A bank's preferred share"], 2),
]),

dict(existing_title="Investment Instruments - Equities, Bonds, and T-Bills", title="Equities, Bonds and T-Bills", video="", content="""
## The three building blocks

### 1. Equities (shares or stocks)
A share is a unit of ownership. Your return comes from two sources:

- **Dividends**: a portion of profits the company chooses to pay to shareholders.
- **Capital gains**: selling the share for more than you paid.

Equities have the **highest potential return and the highest risk**. If the company does badly, the share price can fall below what you paid, and dividends can be skipped.

*Example:* Wegagen Bank (WGBX), Awash Bank (AWAB), Gadaa Bank (GDAB), Abay Bank (ABAYB), Bank of Abyssinia (BOAX), Sidama Bank (SIDAX) and Ethio Telecom (TELE) are all equities listed on the ESX.

### 2. Bonds
A bond is a **loan you make to a government or company**. The issuer promises to pay a fixed **coupon** (interest) at regular intervals and to repay the **face value** on the **maturity date**.

- Lower risk than shares, because the payments are contractual.
- Still has risks: the issuer may fail to pay (**credit risk**), and if market interest rates rise, the market value of your bond falls (**interest-rate risk**).

### 3. Treasury bills (T-bills)
T-bills are **short-term government debt**. In Ethiopia they are auctioned by the National Bank of Ethiopia (NBE) in maturities such as **28, 91, 182 and 364 days**. They are sold at a **discount** to face value: you pay less today and receive the full face value at maturity, and the difference is your interest.

*Illustration:* a 91-day bill with a face value of ETB 100,000 bought at ETB 98,500 earns ETB 1,500 over 91 days, which is about **1.52%** for the period, or about **6.1% a year** (1.52% x 365/91). Always compare returns on an **annualised** basis.

## Risk and return at a glance

| Instrument | Typical risk | Typical return source | Typical horizon |
|---|---|---|---|
| T-bills | Lowest (government-backed) | Discount to face value | Under 1 year |
| Bonds | Low to medium | Coupons + repayment | 1 to 30 years |
| Shares | Medium to high | Dividends + price growth | Long term (5+ years) |

## Key takeaways

1. Shares = ownership (dividends + capital gains); bonds and T-bills = loans (interest).
2. Higher expected return always comes with higher risk.
3. A T-bill earns its interest through a **discount**; always annualise to compare.
""", quiz=[
 ("How does a treasury bill pay you interest?", ["Monthly coupons", "It is sold at a discount and repaid at face value", "Through dividends", "Through bonus shares"], 1),
 ("A 91-day T-bill with face value ETB 100,000 is bought for ETB 98,500. What is the interest earned?", ["ETB 98,500", "ETB 1,500", "ETB 15,000", "ETB 150"], 1),
 ("Which instrument typically carries the highest risk and highest potential return?", ["T-bills", "Bonds", "Equities", "Savings deposits"], 2),
 ("What is the 'coupon' on a bond?", ["The bond's face value", "The periodic interest payment", "The maturity date", "The broker's fee"], 1),
 ("If market interest rates rise, what generally happens to the market price of an existing fixed-rate bond?", ["It rises", "It falls", "It is unaffected", "It is cancelled"], 1),
]),

dict(title="The Players: Who Does What in Ethiopia's Market", video="", content="""
## An ecosystem, not just an exchange

A functioning capital market needs several kinds of participants, each with a different job. Here is Ethiopia's cast as of late 2026.

### The regulator: ECMA
The **Ethiopian Capital Market Authority (ECMA)** was created under the Capital Market Proclamation (No. 1248/2021). It licenses market participants, approves prospectuses, supervises the exchange and protects investors. Think of ECMA as the **referee**. Its office is in the Minaye Corporate Building, Flamingo area, Addis Ababa.

### The marketplace: ESX
The **Ethiopian Securities Exchange (ESX)** is the venue where buyers and sellers are matched. It began operations with the listing of Wegagen Bank on **10 January 2025**. ESX does not give you investment advice or take your orders directly; you reach it **through licensed intermediaries**.

### The intermediaries: capital market service providers (CMSPs)
ECMA licenses firms in categories that you can browse in Stocet's **CMSPs** directory. As of October 2026 there are 20 licensees:

- **Investment banks (7)**: arrange and underwrite share and bond offerings, and advise companies that want to list. Examples: CBE Capital, Awash Capital, Wegagen Capital, First Addis, Prime Capital, Siinqee and PrizeWorth.
- **Securities dealers (2)**: the firms through which investors **buy and sell** on the ESX: Ethio-Fidelity Securities and Gadaa Securities Dealer.
- **Securities investment advisors (11)**: advise companies and investors on capital-market transactions.

### Issuers
Companies (and in future governments and agencies) that sell securities to raise money. Today the seven ESX-listed issuers are six banks and Ethio Telecom.

### Investors
Retail investors (individuals) and institutions (banks, insurers, funds). Your rights as an investor include receiving timely, accurate company disclosures, voting at shareholder meetings (AGMs) and receiving your share of dividends declared.

### The central bank: NBE
The **National Bank of Ethiopia** is not a capital market regulator, but it matters enormously: it sets monetary policy (interest rates), runs the T-bill auctions, supervises the banks that make up most ESX-listed companies, and manages foreign exchange policy.

## Key takeaways

1. ECMA = regulator (the referee); ESX = marketplace; CMSPs = licensed intermediaries.
2. You trade through a licensed **securities dealer**, not directly with the exchange.
3. Always check that a firm is on ECMA's licensee list before giving it money.
""", quiz=[
 ("Which body licenses capital market participants and approves prospectuses in Ethiopia?", ["ESX", "ECMA", "Ministry of Finance", "Commercial Bank of Ethiopia"], 1),
 ("Through whom do retail investors place orders to buy ESX-listed shares?", ["Directly with the NBE", "A licensed securities dealer", "Any bank branch", "The company's CEO"], 1),
 ("On what date did the ESX begin operations with its first listing?", ["10 January 2025", "1 July 2024", "13 June 2025", "28 July 2026"], 0),
 ("What do investment banks do in the market?", ["Run T-bill auctions", "Arrange and underwrite offerings and advise issuers", "Set interest rates", "Audit company accounts"], 1),
 ("Before giving money to a firm that offers 'ESX investing', what should you do first?", ["Check it appears on ECMA's licensee list", "Ask for a guaranteed return", "Pay a deposit to hold a spot", "Nothing; all firms are licensed"], 0),
]),

dict(title="Risk, Return and Diversification", video="", content="""
## No free lunch

**Return** is what you earn; **risk** is the chance that the outcome is different from what you expected, usually worse. The two are linked: **to earn higher expected returns, you must accept more uncertainty.**

## Types of risk you will meet on the ESX

- **Market risk**: the whole market falls, taking your shares with it.
- **Company-specific risk**: something bad happens to one company (a loss, a scandal, a change of management).
- **Liquidity risk**: you want to sell, but there are not enough buyers at a fair price. The ESX is young and many shares trade thinly, so this is a **real and important** risk in Ethiopia.
- **Inflation risk**: your return is lower than inflation, so your money buys less. With Ethiopian inflation in the double digits (13.4% in May 2026), a 10% return is a **negative real return**.
- **Concentration risk**: too much of your money in one company or sector. All six ESX-listed banks react to the same things (NBE policy, credit rules, the economy), so owning three banks is **not** real diversification.

## Real return

Real return is approximately **nominal return minus inflation**. If a T-bill pays 9% and inflation is 13%, your purchasing power falls by about 4% a year. Always ask: *what is the return after inflation?*

## Diversification

Diversification means spreading money across investments that **do not all move together**, so one bad outcome does not sink you. Ways to diversify:

1. **Across companies** (not just one stock).
2. **Across sectors** (banks and telecom, and as new listings arrive, insurance, manufacturing, agriculture).
3. **Across asset types** (shares, bonds, T-bills, cash).
4. **Across time**: invest regularly rather than all at once (called pound-cost or *birr-cost averaging*).

Diversification **reduces company-specific risk but cannot remove market risk**.

## Know your own risk profile

Ask yourself: *If this investment halved in value next year, would I be able to hold on, or would I need the money?* Money you need within 1 to 3 years does not belong in shares.

## Key takeaways

1. Higher expected return comes with higher risk.
2. On the ESX, watch **liquidity**, **concentration** and **inflation** risk especially.
3. Diversify across companies, sectors and asset types; it cuts company risk, not market risk.
4. Only invest in shares money you will not need soon.
""", quiz=[
 ("Inflation is 13% and your investment returns 9% a year. Roughly what is your real return?", ["+22%", "+4%", "-4%", "0%"], 2),
 ("Why is owning shares in three different ESX-listed banks only weak diversification?", ["Banks never pay dividends", "They all react to the same economic and regulatory factors", "Banks cannot be sold", "Because shares are illiquid"], 1),
 ("What is liquidity risk?", ["The risk a company runs out of cash", "The risk you cannot sell quickly at a fair price", "The risk of inflation", "The risk of fraud"], 1),
 ("Diversification mainly reduces which kind of risk?", ["Market risk", "Inflation risk", "Company-specific risk", "Regulatory risk"], 2),
 ("Which money is most appropriate to put into shares?", ["Rent due next month", "School fees due in 6 months", "Savings you will not need for 5+ years", "An emergency fund"], 2),
]),

dict(title="The Time Value of Money and Compounding", video="", content="""
## A birr today is worth more than a birr tomorrow

If someone offers you ETB 10,000 today or ETB 10,000 in a year, you take it today, because you can invest it, and inflation will erode the future amount. This is the **time value of money**, the foundation of nearly all of finance.

## Future value

If you invest an amount **P** at an annual return **r** for **n** years:

**Future value = P x (1 + r)^n**

*Example:* ETB 100,000 at 15% a year for 5 years = 100,000 x 1.15^5 = **ETB 201,136**. Your money doubled in about 5 years.

## Compounding

Compounding means **earning returns on your returns**. It starts slowly and then accelerates:

| Year | Value at 15% a year |
|---|---|
| 0 | 100,000 |
| 5 | 201,136 |
| 10 | 404,556 |
| 15 | 813,706 |
| 20 | 1,636,654 |

The same ETB 100,000 grows 16-fold in 20 years, not 4-fold (which simple, non-compounded 15% would give: 100,000 + 20 x 15,000 = 400,000). **Time is the investor's biggest advantage**, so starting early matters more than starting big.

## The Rule of 72

A quick estimate: **years to double = 72 / annual return (%)**. At 12% a year, money doubles in about 6 years; at 6%, about 12 years.

## Present value and discounting

The reverse question: *what is a future amount worth today?* **Present value = Future value / (1 + r)^n**. ETB 200,000 received in 5 years, discounted at 15%, is worth about ETB 99,000 today. Valuing a share by estimating its future dividends and discounting them is called the **dividend discount model**, and it is the basis of the valuation methods in the Fundamental Analysis course.

## Nominal vs. real

Always check whether a return figure is before or after inflation. A 15% nominal return with 13% inflation is only about **1.8% real** (1.15 / 1.13 - 1).

## Key takeaways

1. Money today is worth more than money later.
2. Future value = P x (1 + r)^n; compounding makes growth accelerate.
3. The Rule of 72 gives a quick doubling time.
4. Compare returns after inflation.
""", quiz=[
 ("What is the future value of ETB 100,000 invested at 15% a year for 5 years (compounded)?", ["ETB 175,000", "ETB 201,136", "ETB 150,000", "ETB 115,000"], 1),
 ("Using the Rule of 72, about how long does money take to double at 12% a year?", ["3 years", "6 years", "12 years", "72 years"], 1),
 ("What does compounding mean?", ["Earning returns on your earlier returns", "Paying fees on every trade", "Investing in several companies", "Receiving a fixed coupon"], 0),
 ("A 15% nominal return with 13% inflation gives a real return of roughly:", ["28%", "15%", "2%", "-2%"], 2),
 ("Which factor is the investor's biggest advantage in compounding?", ["Trading frequently", "Time", "Borrowing", "Brokers' advice"], 1),
]),
]
