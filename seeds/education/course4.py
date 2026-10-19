COURSE = dict(
    key="c4",
    title="Technical Analysis & the Stocet Trading Workspace",
    duration="2 hours 30 mins",
    description="Learn to read price charts and use indicators, then put the ideas to work in Stocet Analytics: candlesticks, trends, support and resistance, moving averages, RSI, MACD, Bollinger Bands, drawing tools, bar replay and backtesting. Includes the honest limits of technical analysis in a young, thin market.",
)

LESSONS = [
dict(title="Reading Price Charts: Candlesticks and Timeframes", video="", content="""
## What a chart shows

A price chart plots price against time. Stocet's **Analytics** workspace draws charts for every ESX ticker using the prices collected from the exchange.

## The candlestick

Each candle summarises one time period (a minute, an hour, a day):

- **Open**: the first price in the period.
- **Close**: the last price in the period.
- **High** and **Low**: the extremes, shown as thin lines called **wicks**.
- The **body** is the range between open and close. It is **green** when the close is above the open (price rose) and **red** when the close is below the open.

Long bodies show strong conviction; long wicks show that price was pushed to an extreme and rejected.

## Other chart types in Stocet

- **Hollow candles** and **OHLC bars**: the same data, different style.
- **Heikin Ashi**: smoothed candles that make trends easier to see.
- **Line / Area**: closing prices only.

## Timeframes

Stocet supports **1m, 5m, 1h and 1D** candles. Rules of thumb:

- **Longer timeframes (daily)** show the real trend and produce fewer false signals.
- **Shorter timeframes** show more detail but more noise.
- Look at the **bigger timeframe first**, then zoom in.

## A caution about ESX data

The ESX is young, and many shares trade a few times a day. That means:

- Intraday candles can have **gaps and flat stretches** where no trade happened.
- A candle formed by **one trade** says very little.
- **Daily** charts are the most meaningful starting point.

Stocet's chart does not currently show volume because the price feed does not include traded volume. Be cautious of any pattern that depends on volume.

## Try it in Stocet

Open **Analytics**, choose a ticker, set the timeframe to **1D** and the range to **1Y**, and switch between *Candles*, *Heikin Ashi* and *Line*. Hover the chart to see open, high, low and close in the legend at the top-left.

## Key takeaways

1. A candle shows open, high, low and close for one period.
2. Start with daily charts; short timeframes are noisy, especially on thin ESX stocks.
3. Volume is not available in Stocet's feed.
""", quiz=[
 ("A green candle means:", ["The close was below the open", "The close was above the open", "There was no trade", "Volume was high"], 1),
 ("What do the thin 'wicks' on a candle show?", ["The open and close", "The highest and lowest prices in the period", "The dividend", "The average price"], 1),
 ("Why are daily charts a better starting point for thinly traded ESX shares?", ["They show volume", "They have fewer noisy, near-empty candles", "They update faster", "They show tax"], 1),
 ("Which timeframes does Stocet support?", ["1m, 5m, 1h, 1D", "Only 1D", "Only 1m", "1s and 1W"], 0),
 ("Why should you be careful with a candle formed by a single trade?", ["It is always wrong", "It carries very little information about demand", "It means the share is delisted", "It is a dividend candle"], 1),
]),

dict(title="Trends, Support and Resistance", video="", content="""
## The trend is your friend (until it ends)

- **Uptrend**: a series of **higher highs and higher lows**.
- **Downtrend**: **lower highs and lower lows**.
- **Range (sideways)**: price bounces between a ceiling and a floor.

Identify the trend first. Buying in an uptrend is usually easier than trying to catch a falling share.

## Support and resistance

- **Support**: a price area where buying has repeatedly stopped a fall. Buyers appear there.
- **Resistance**: a price area where selling has repeatedly stopped a rise.

Why they work: traders remember those prices. Someone who bought at 1,000 and watched it fall may sell when it returns to 1,000 ("get out even"), creating resistance.

Key ideas:

1. **The more times a level has held, the more it matters** (but the more likely it is eventually to break).
2. When resistance breaks, it often becomes **support**, and the reverse.
3. Levels are **zones**, not exact numbers. Draw them through clusters of highs or lows.
4. A **breakout** is more convincing if price **closes** beyond the level rather than just spiking through.

## Trendlines

A **trendline** connects at least two swing lows in an uptrend (or two swing highs in a downtrend). The more touches, the more significant. A close beyond the line suggests the trend may be changing.

## Fibonacci retracements

After a move up, price often pulls back part of it before continuing. Traders watch the **23.6%, 38.2%, 50% and 61.8%** retracement levels as possible support. They are not magic: they are useful because many people watch them, and as *reference zones*, not guarantees.

## Try it in Stocet

In **Analytics**, use the drawing tools on the left of the chart: **Trend line**, **Horizontal line**, **Rectangle** for a support zone, **Fib retracement**, and **Measure** to see the percent move and bar count between two points. Drawings are saved for each ticker; click one to select, drag its handles to adjust, and press **Delete** to remove it. Save your whole setup with **Layouts**.

## Key takeaways

1. Uptrend = higher highs and lows; downtrend = lower highs and lows.
2. Support and resistance are zones where buyers or sellers have repeatedly stepped in.
3. Broken resistance often becomes support.
4. Closing beyond a level is stronger evidence than a brief spike.
""", quiz=[
 ("Which describes an uptrend?", ["Lower highs and lower lows", "Higher highs and higher lows", "No clear highs or lows", "A single large candle"], 1),
 ("What is 'support'?", ["A price area where selling has repeatedly stopped rises", "A price area where buying has repeatedly stopped falls", "The dividend level", "The all-time high"], 1),
 ("A broken resistance level often becomes:", ["Irrelevant", "Support", "The next resistance", "The stop-loss"], 1),
 ("A breakout is more convincing when price:", ["Briefly spikes through the level", "Closes beyond the level", "Opens at the level", "Trades on low volume"], 1),
 ("Which Stocet tool shows the percent move and number of bars between two points?", ["Ray", "Measure", "Note", "Rectangle"], 1),
]),

dict(title="Moving Averages and MACD: Following the Trend", video="", content="""
## Moving averages (MA)

A moving average smooths price by averaging the last N closes, so you see direction instead of noise.

- **SMA (simple moving average)**: the plain average of the last N closes.
- **EMA (exponential moving average)**: gives more weight to recent prices, so it reacts faster.
- **HMA, WMA**: faster variants available in Stocet.

Common lengths: **20** (short term), **50** (medium), **200** (long term).

How to read them:

1. **Slope and position**: price above a rising MA suggests an uptrend; below a falling MA, a downtrend.
2. **Crossovers**: when a fast MA crosses above a slow MA (e.g. EMA 9 over EMA 21), it is a bullish signal, and the reverse is bearish. These are *lagging*: they confirm a move after it has started.
3. **Dynamic support/resistance**: in strong trends, price often bounces off the 20 or 50 MA.

**Weakness:** in a sideways market, MAs whip back and forth and give many **false signals**.

## MACD

**MACD (Moving Average Convergence Divergence)** = EMA(12) - EMA(26). A 9-period EMA of the MACD is the **signal line**, and the difference between them is the **histogram**.

- MACD crossing **above** the signal line: bullish momentum.
- MACD crossing **below** the signal line: bearish momentum.
- MACD above zero: the short-term trend is above the long-term trend.
- **Divergence**: price makes a new high but MACD does not; momentum may be fading.

## Supertrend and Parabolic SAR

Two trend-following tools in Stocet that trail price and flip when the trend turns; they are popular for placing trailing stops.

## Try it in Stocet

Click **Indicators**, search for *EMA*, *MACD* or *Supertrend*, and add them. Open the gear in the chart legend to change the lengths and watch the chart update instantly. Then open the **Strategy** tab, pick **EMA crossover** or **MACD signal cross**, and the backtester shows every buy and sell it would have made on this chart.

## Key takeaways

1. Moving averages smooth price; EMAs react faster than SMAs.
2. Crossovers lag; they work best in trending markets and badly in ranges.
3. MACD shows momentum; watch signal-line crosses and divergence.
""", quiz=[
 ("Which moving average reacts faster to recent price changes?", ["SMA", "EMA", "They are identical", "None"], 1),
 ("In which market condition do moving-average crossovers perform worst?", ["Strong uptrend", "Strong downtrend", "Sideways range", "After a breakout"], 2),
 ("MACD is calculated as:", ["EMA(12) minus EMA(26)", "SMA(50) minus SMA(200)", "Price minus EMA(9)", "High minus low"], 0),
 ("MACD crossing above its signal line is generally read as:", ["Bearish momentum", "Bullish momentum", "A dividend", "No information"], 1),
 ("Why are MA crossovers called 'lagging' signals?", ["They predict the future", "They confirm a move after it has already begun", "They use volume", "They are only for banks"], 1),
]),

dict(title="Momentum and Volatility: RSI, Stochastic, Bollinger Bands, ATR", video="", content="""
## Momentum oscillators

Oscillators move between fixed bounds and show how stretched price is.

### RSI (Relative Strength Index)
Scale **0 to 100**, default length 14.

- **Above 70**: overbought: price rose fast and may be due a pause.
- **Below 30**: oversold: price fell fast and may be due a bounce.
- In strong trends RSI can stay overbought or oversold for a long time. *Overbought is not a sell signal by itself.*
- **Divergence** (price makes a new low, RSI does not) often precedes reversals.

### Stochastic and Williams %R
They show where the close sits within the recent high-low range. Similar idea to RSI, faster and noisier.

### CCI and ROC
Measure how far price is from its average, or the percentage change over N bars.

## Volatility tools

### Bollinger Bands
A 20-period moving average with bands **2 standard deviations** above and below.

- **Squeeze**: bands narrow = low volatility; a big move often follows (direction unknown).
- Price touching the upper band in an uptrend can be strength, not a sell signal.
- In a range, the bands help find mean-reversion trades.

### ATR (Average True Range)
The average size of recent bars, in price terms. ATR does not give direction; it measures **how much the share typically moves**. That makes it the best tool for **setting stop-losses and position sizes**: a stop of 2 x ATR sits outside normal noise.

### ADX
Measures **trend strength** (not direction). Above about 25 suggests a trending market where trend tools work; below 20, a range where mean-reversion tools work better.

## Match the tool to the market

| Market | Better tools |
|---|---|
| Trending | Moving averages, MACD, Supertrend, Donchian breakouts |
| Range-bound | RSI, Stochastic, Bollinger reversion |

## Try it in Stocet

Add **RSI** and **Bollinger Bands**, then in the Strategy tab choose **RSI mean reversion** and **Bollinger reversion** and compare their results on a range-bound stock vs. a trending one.

## Key takeaways

1. RSI above 70 or below 30 shows stretch, not a guaranteed reversal.
2. Bollinger squeezes warn of coming volatility; ATR measures typical movement and drives stops.
3. Use trend tools in trends and mean-reversion tools in ranges; ADX helps decide which.
""", quiz=[
 ("An RSI reading above 70 is usually called:", ["Oversold", "Overbought", "Neutral", "A breakout"], 1),
 ("Which indicator is best suited to setting stop-loss distances?", ["ATR", "RSI", "MACD histogram", "Stochastic"], 0),
 ("A Bollinger Band 'squeeze' suggests:", ["Volatility is low and a bigger move may follow", "Price will certainly fall", "A dividend is due", "The market is closed"], 0),
 ("ADX measures:", ["Trend direction", "Trend strength", "Volume", "Price level"], 1),
 ("In a strong uptrend, RSI above 70:", ["Always means sell now", "Can persist for a long time", "Means the data is wrong", "Cannot happen"], 1),
]),

dict(title="Backtesting and Replay: Testing Ideas Without Risk", video="", content="""
## Why test?

Before risking money on a rule like "buy when the 9 EMA crosses above the 21 EMA", ask: *how would it have performed?* That is a **backtest**: applying the rule to historical data.

Stocet's **Strategy** tab runs backtests on the chart you are viewing:

- Choose a **preset** (SMA/EMA crossover, RSI reversion, MACD, Bollinger, Donchian, Stochastic, Price vs MA, Buy & hold) or build **custom rules** (for example: "Close crosses above EMA 50 AND RSI is above 50").
- Set **risk management**: stop loss, take profit, trailing stop, breakeven, time stop and position sizing.
- Set realistic **costs**: commission and slippage.
- Read the **Performance** panel: net profit, return vs. buy-and-hold, maximum drawdown, win rate, profit factor, expectancy, Sharpe and Sortino, and the monthly returns table. The **Trades** tab lists every trade with how it ended.

## How to read the results

- **Net profit vs. buy & hold.** A strategy that earns 12% when simply holding earned 35% did *worse* despite being profitable.
- **Maximum drawdown**: the worst peak-to-trough loss. Could you have sat through it?
- **Win rate** is less important than **payoff**: a 40% win rate can be profitable if winners are much larger than losers.
- **Profit factor** = gross profit / gross loss. Above 1.5 is decent.
- **Number of trades**: with fewer than about 20, the result is mostly luck.
- **Expectancy**: the average profit per trade after costs.

## The traps

1. **Overfitting.** If you keep tweaking parameters until the past looks perfect, you are fitting noise. Stocet's **Optimize** tab ranks settings on the first 70% of the chart and shows results on the last 30%; if the "best" setting falls apart out-of-sample, it was overfit.
2. **Too little data.** The ESX has less than two years of history. Results from a few months say little.
3. **Ignoring costs.** Set commission to your dealer's real rate. See how results change.
4. **Look-ahead bias.** Stocet fills signals at the *next* bar's open to avoid this; be suspicious of any tool that fills at the signal bar's close.
5. **Illiquidity.** The backtest assumes you can always trade at the open price. In thin shares you may not.
6. **Past is not future.** A regime (such as high inflation and rapid bank growth) can change.

## Bar replay: practice without risk

Press **Replay** to hide the future and step through history bar by bar. Place **paper trades** with stop-loss, take-profit and trailing stops, and watch them play out. Use it to practise the *process* (entering, sizing, managing, exiting) before using real money, and save the trades to your **Journal**.

## Key takeaways

1. Backtests tell you whether an idea *would have* worked, not whether it will.
2. Compare against buy-and-hold, check drawdown and trade count, and include costs.
3. Use out-of-sample checks to avoid overfitting; use Replay to practise.
""", quiz=[
 ("A strategy returned 12% while buy-and-hold returned 35%. What does this mean?", ["The strategy won", "The strategy was profitable but did worse than simply holding", "The data is wrong", "The strategy is risk-free"], 1),
 ("What is overfitting?", ["Tuning rules so tightly to past data that they only fit noise", "Using too many indicators on screen", "Trading too small", "Paying too much commission"], 0),
 ("Why does Stocet fill signals at the next bar's open?", ["To save time", "To avoid look-ahead bias", "To increase profits", "Because exchanges only open once"], 1),
 ("What does Stocet's Optimize tab's out-of-sample result tell you?", ["How the setting did on data not used to rank it", "The best possible future return", "The commission", "The dividend yield"], 0),
 ("What is bar replay mainly for?", ["Predicting prices", "Practising decisions and risk management on past data without real money", "Placing real orders", "Downloading data"], 1),
]),
]
