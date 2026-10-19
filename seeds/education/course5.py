COURSE = dict(
    key="c5",
    title="Risk Management & Investor Psychology",
    duration="2 hours",
    description="The skills that decide whether you keep your money: position sizing with the 1% rule, stop-loss, take-profit and trailing stops, reward-to-risk, portfolio construction, behavioural biases, and keeping a trade journal. Includes hands-on exercises with Stocet's trade planner and replay simulator.",
)

LESSONS = [
dict(title="Position Sizing: How Much Should You Buy?", video="", content="""
## Survive first, profit second

Most beginners ask "what should I buy?". Professionals ask first: **"how much can I lose on this trade, and is that acceptable?"** Position sizing answers that.

## The fixed-risk rule

Decide the **maximum percentage of your account you will lose** if a trade goes wrong. A common guideline is **1% per trade** (some use 0.5% to 2%).

**Shares = (Account size x Risk %) / (Entry price - Stop price)**

*Worked example:* Account ETB 200,000. Risk 1% = ETB 2,000. You plan to buy at ETB 1,000 with a stop at ETB 950, so the risk per share is ETB 50.

Shares = 2,000 / 50 = **40 shares** (position value ETB 40,000, 20% of the account).

If the stop is hit, you lose ETB 2,000 (plus costs), which is 1% of the account. Place the stop wider at ETB 900 and the same rule gives only 20 shares. **The stop distance decides the size**, and not your enthusiasm.

## Why small risk matters: the maths of recovery

| Loss | Gain needed to get back to even |
|---|---|
| -10% | +11% |
| -25% | +33% |
| -50% | +100% |
| -75% | +300% |

Large losses are brutal. With 1% risk per trade you could lose ten trades in a row and still be down only about 10%.

## Practical limits

1. **Cash cap**: you cannot buy more than you have; the formula may suggest more shares than you can afford.
2. **Liquidity cap**: on a thin ESX share you may not be able to sell a large position quickly. Keep positions modest relative to normal daily trading.
3. **Concentration cap**: even at 1% risk, do not put more than, say, 20% to 25% of the account in one company.
4. **Costs**: include commission, which makes the real risk slightly larger.

## Try it in Stocet

Open **Analytics → Trade**. Set your **account size** and **risk per trade**, click **New plan**, drag the stop line on the chart, and press **size by risk**. The panel shows your maximum loss, maximum gain and reward-to-risk before you commit.

## Key takeaways

1. Size positions from the **risk you accept**, not the amount you hope to make.
2. Shares = (account x risk %) / (entry - stop).
3. Small, consistent risk keeps you in the game; big losses are hard to recover from.
""", quiz=[
 ("Account ETB 300,000, risk 1%, entry ETB 600, stop ETB 570. How many shares?", ["50", "100", "200", "500"], 1),
 ("If you widen the stop-loss and keep 1% risk, the position size should:", ["Increase", "Decrease", "Stay the same", "Double"], 1),
 ("After a 50% loss, what gain is needed to get back to even?", ["25%", "50%", "100%", "150%"], 2),
 ("Why add a liquidity cap in ESX shares?", ["Taxes", "You may not be able to exit a large position quickly in a thin market", "Dividends are lower", "It is required by law"], 1),
 ("What decides your position size in the fixed-risk method?", ["How confident you feel", "The distance from entry to stop", "The company's logo", "Last week's profit"], 1),
]),

dict(title="Stop-Loss, Take-Profit and Trailing Stops", video="", content="""
## Planning the exit before the entry

Every trade needs three prices decided **before** you buy:

1. **Entry**: where you buy.
2. **Stop-loss**: where you admit you were wrong and sell.
3. **Take-profit (target)**: where you plan to sell for a gain.

## Stop-loss

A stop-loss caps your loss. Place it where the trade idea is **proven wrong**, not at a round number:

- **Below support** or a recent swing low.
- **ATR-based**: entry minus 1.5 to 3 times ATR (keeps you outside normal noise).
- **Percentage-based**: simple (e.g. 5% below entry) but ignores the share's own volatility.

Never move a stop **further away** to avoid taking a loss.

> **Important on the ESX:** *Stocet cannot place stops for you*, and your orders are placed on ESX Neway with your licensed dealer. If your dealer does not support stop orders, a stop is only a **plan**: you must enter the sell yourself when price reaches it. Use Stocet's **price alerts** (created from the Execute screen) so you are notified. Also remember that if a share is halted or gaps through your stop, you may sell at a worse price, and in thin markets there may be no buyers at your stop.

## Take-profit and reward-to-risk

**Reward-to-risk (R:R)** = (target - entry) / (entry - stop).

- Entry 1,000, stop 950, target 1,150: risk 50, reward 150 -> **R:R = 3**.
- With R:R = 2, you only need to win **34%** of trades to break even; with R:R = 1, you need 50%.

The formula: **break-even win rate = 1 / (1 + R:R)**.

Look for setups with R:R of **2 or more**. Stocet's trade planner displays R:R and the break-even win rate as you drag the lines.

## Trailing stops

A **trailing stop** moves up as price rises and never moves down, locking in profit:

- **Percentage trail**: stop stays N% below the highest price since entry.
- **ATR trail**: stop stays N x ATR below the high.
- **Breakeven stop**: once the trade is up by some amount, move the stop to your entry so the trade can no longer lose.

Trailing stops let winners run but can exit early on a normal pullback. Backtest different trail sizes in Stocet's **Strategy** tab and compare.

## Time stops

If a trade has not worked after N bars, exit and free the capital. Stocet's backtester has a **time stop** setting.

## Try it in Stocet (Replay)

Start **Replay**, make a plan with entry, stop and target, place a market buy, and press play. Watch whether the stop, target or trailing stop is hit. Do this 20 times and you will learn more than from any article.

## Key takeaways

1. Decide entry, stop and target before you trade.
2. Put stops where the idea is wrong; never widen them to avoid a loss.
3. Aim for R:R of 2 or more; break-even win rate = 1 / (1 + R:R).
4. On the ESX, stops are often manual: use alerts and know your dealer's order types.
""", quiz=[
 ("Entry 1,000, stop 950, target 1,150. What is the reward-to-risk ratio?", ["1", "2", "3", "0.33"], 2),
 ("With a reward-to-risk of 2, roughly what win rate do you need to break even?", ["20%", "34%", "50%", "66%"], 1),
 ("What does a trailing stop do?", ["Moves up as price rises and never moves down", "Moves randomly", "Is the same as a limit order", "Only works after market close"], 0),
 ("Why can a stop-loss fail to protect you fully in a thin market?", ["Stops never work", "Price may gap through it or there may be no buyers at your level", "Stops cost too much", "ESX forbids stops"], 1),
 ("Which is a poor use of stops?", ["Placing it below support", "Using ATR for distance", "Moving it further away to avoid taking the loss", "Moving it to breakeven after a gain"], 2),
]),

dict(title="Building a Portfolio: Diversification and Rebalancing", video="", content="""
## From single trades to a portfolio

Your results come from the **whole portfolio**, not one trade. A good portfolio is built from your goals, time horizon and risk tolerance.

## Step 1: Goals and horizon

- **Emergency fund first**: 3 to 6 months of expenses in safe, accessible savings (bank deposits or T-bills) before investing in shares.
- **Short goals (under 2 years)**: safe assets (T-bills, deposits).
- **Medium (2 to 5 years)**: a mix.
- **Long (5+ years)**: more shares.

## Step 2: Asset allocation

How you divide money between **asset classes** matters more than which stock you pick. A sample, not advice:

| Profile | Shares | Bonds / T-bills | Cash |
|---|---|---|---|
| Cautious | 20% | 60% | 20% |
| Balanced | 50% | 40% | 10% |
| Growth | 70% | 25% | 5% |

Stocet's **Advisory** tool can turn a short questionnaire into a suggested allocation (it is informational, not a licensed advisory service).

## Step 3: Spread within shares

With today's ESX, you can own up to seven companies, six of them banks. Practical implications:

- Holding several banks is **not** true diversification; they share the same drivers.
- Consider **Ethio Telecom (TELE)** as the non-bank exposure, and add others as they list.
- Assets outside the ESX (T-bills, deposits, property, a business) matter for the other half of your risk.

## Step 4: Monitor and rebalance

If shares do well, they become a larger share of your portfolio than planned. **Rebalancing** means trimming what grew and topping up what lagged, to bring risk back to target. Do it on a schedule (once or twice a year) or when allocations drift by 5 to 10 percentage points.

## Step 5: Track performance honestly

Use Stocet's **Portfolio** page to record holdings with cost basis and see profit and loss, allocation by sector and concentration. Compare against a fair benchmark: **buy-and-hold of the ESX stocks you own**, **T-bills**, and **inflation**.

## Common mistakes

- Putting everything into one hot stock.
- Treating a few weeks of gains as skill.
- Not having cash for opportunities or emergencies.
- Ignoring inflation when comparing returns.

## Key takeaways

1. Emergency fund first, then goals, then allocation.
2. Several banks do not equal diversification; look outside the ESX too.
3. Rebalance periodically and compare against T-bills and inflation.
""", quiz=[
 ("What should you have before investing in shares?", ["A hot tip", "An emergency fund", "A margin loan", "Three brokers"], 1),
 ("Why is owning several ESX-listed banks limited diversification?", ["Banks are illiquid", "They all depend on the same economic and regulatory drivers", "Banks do not pay dividends", "It is illegal"], 1),
 ("What is rebalancing?", ["Selling everything", "Bringing allocations back to target after they drift", "Changing brokers", "Switching to a new currency"], 1),
 ("Money needed in under two years should mostly be in:", ["Shares", "Safe assets like T-bills and deposits", "A single stock", "Pre-IPO deals"], 1),
 ("A fair benchmark for your portfolio includes:", ["Your friend's returns", "T-bills and inflation", "Social-media claims", "Only your best month"], 1),
]),

dict(title="The Psychology of Investing: Biases and Discipline", video="", content="""
## You are your biggest risk

Markets are a test of behaviour. The same mistakes recur everywhere, and they are well documented.

## Common biases

1. **Loss aversion.** A loss hurts about twice as much as an equal gain feels good. Result: people **hold losers too long** (hoping to get back to even) and **sell winners too early** (to lock in the pleasure).
2. **Overconfidence.** After a few wins we credit skill and trade larger and more often. Most of the early success in a rising market is the market.
3. **Confirmation bias.** We seek information that supports what we already believe and ignore the rest.
4. **Herd behaviour / FOMO.** Buying because everyone else is, after the price has already risen.
5. **Anchoring.** Fixating on a price we saw (our entry, or last month's high) and judging everything against it.
6. **Recency bias.** Assuming the recent past will continue: a stock that doubled "will keep going", one that fell "will keep falling".
7. **Sunk-cost fallacy.** "I have already lost so much; I cannot sell now." The market does not know or care what you paid.
8. **Gambler's fallacy and revenge trading.** Trying to win back a loss by trading bigger.

## In a young market

The ESX is new and exciting. Social media and group chats amplify FOMO, and thin trading means a few buyers can push a price up temporarily. Be especially alert to herd behaviour and rumours.

## The antidotes

- **A written plan** for each trade: why, entry, stop, target, size. (Stocet's Trade Plan does this.)
- **Rules over feelings**: if the stop is hit, you sell, with no negotiation.
- **Pre-commit to your risk** (1% per trade) so a bad day cannot become a disaster.
- **A cooling-off rule**: after a loss of a set size, stop for the day.
- **Review, do not react**: check the portfolio on a schedule, not every ten minutes.
- **Seek disconfirming evidence**: before buying, write down the best argument *against* the trade.

## The trade journal

Keep a record for every trade: date, ticker, reason, entry/stop/target, size, result, and **how you felt**. Review monthly. Patterns appear: "I always exit winners early on Mondays", "my trades after a loss are worse". Stocet's **Journal** tab stores planned, replay and real trades, with profit and loss and win rate.

## Key takeaways

1. Loss aversion makes us hold losers and sell winners; rules fix that.
2. Write a plan before each trade and follow it.
3. Keep a journal and review it; it is the cheapest investment education.
""", quiz=[
 ("Loss aversion typically makes investors:", ["Sell losers quickly and hold winners", "Hold losers too long and sell winners too early", "Trade less", "Ignore costs"], 1),
 ("Buying a share mainly because everyone else is, after its price has risen, is:", ["Anchoring", "Herd behaviour / FOMO", "Diversification", "Rebalancing"], 1),
 ("'I have already lost so much, I cannot sell now' is an example of:", ["Sunk-cost fallacy", "Confirmation bias", "Recency bias", "Diversification"], 0),
 ("What is a useful antidote to emotional trading?", ["Checking prices every minute", "A written plan with rules for entry, stop and target", "Doubling the size after a loss", "Following group chats"], 1),
 ("What should a trade journal include?", ["Only winners", "Reason, entry/stop/target, size, result and how you felt", "Only the profit", "Only the ticker"], 1),
]),
]
