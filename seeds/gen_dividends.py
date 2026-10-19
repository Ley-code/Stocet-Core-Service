#!/usr/bin/env python3
"""Generates seeds/004_dividends.sql: real dividend and capital-raising events for ESX-listed companies.

Convention: corporate_actions.amount = ETB per share (derived = total declared / shares in issue at the profit-year end, shown in the
description); NULL when the per-share figure cannot be derived from published data. event_date is the exact AGM date when an AGM notice
exists, otherwise the fiscal year-end of the profit year the dividend is paid out of.
Sources: ESX AGM notices (BOA 12 Oct 2026, Wegagen 21 Sep 2026, Gadaa 21 Sep 2026), audited FY2025/26 statements, Awash Annual Report 2024/25,
Sidama Annual Report 2024/25.
"""
def q(s): return "$q$" + s + "$q$"
def n(v): return "NULL" if v is None else repr(round(v, 4))

ESX = "https://esx.et/wp-content/uploads"
A = []  # ticker, type, title, event_date, description, amount

def add(t, typ, title, date, desc, amount=None): A.append((t, typ, title, date, desc, amount))

# ---- Bank of Abyssinia (BOAX)
add("BOAX", "dividend", "Dividend approved at 30th AGM, on FY2025/26 profit", "2026-10-10",
    "The AGM of 10 Oct 2026 approved a dividend of ETB 6,033,305,586.09, which equals the bank's whole FY2025/26 retained earnings and about 52.3% of net profit (ETB 11.54bn). Per-share figure is approximate: total / 17,481,120 shares in issue. "
    f"Source: BOA notification to ESX dated 12 Oct 2026 ({ESX}/2026/10/Notification-of-17th-AGM-Resolution-BOA.pdf).", 6033305586.09 / 17481120)
add("BOAX", "agm", "30th Annual Ordinary & 17th Extraordinary General Meetings", "2026-10-10",
    "Elected eight of nine existing directors for a further three-year term plus four new independent directors (subject to NBE approval); approved the dividend; approved raising paid-up capital by ETB 10bn through a rights and public offering, without prejudice to pre-emptive rights.")
add("BOAX", "rights_issue", "Rights offer allotted: 2,481,120 shares at ETB 1,600", "2026-06-24",
    "Rights offer under the prospectus approved 29 Jan 2026. ECMA approved the allotment of 2,481,120 shares worth ETB 3,969,792,000 on 24 Jun 2026 (par ETB 1,000 + ETB 600 premium); ETB 2.66bn of the amount was paid by shareholders from their dividend and ETB 1.31bn in cash. Paid-up capital rose from ETB 15bn to ETB 17,481,120,000 on 26 Jun 2026. Source: audited FY2025/26 statements, note 28.", 1600.0)
add("BOAX", "rights_issue", "EGM approves ETB 10bn capital increase (rights and public offering)", "2026-10-10",
    "The Extraordinary General Meeting approved a proposal to increase the bank's paid-up capital by ETB 10 billion through a rights and public offering, without prejudice to shareholders' pre-emptive rights. Terms and timing not yet published.")
add("BOAX", "dividend", "Dividend declared on FY2024/25 profit (recorded in FY2025/26)", "2025-06-30",
    "ETB 4,973,379,000 declared (about 68% of FY2024/25 profit of ETB 7.30bn); ETB 2,463,722,000 was paid in cash and the rest was reinvested by shareholders in the June 2026 rights offer. Per share on the 15,000,000 shares in issue at 30 Jun 2025. Source: audited FY2025/26 statements.", 4973379000 / 15000000)
add("BOAX", "dividend", "Dividend declared on FY2023/24 profit (recorded in FY2024/25)", "2024-06-30",
    "ETB 2,576,914,000 declared; ETB 1,806,612,000 paid in cash during FY2024/25. Source: audited FY2025/26 statements (comparatives).")

# ---- Wegagen Bank (WGBX)
add("WGBX", "dividend", "Dividend declared at 33rd AGM, on FY2025/26 profit", "2026-09-19",
    "The bank declared a dividend of ETB 2,141,446,943.70, about 68% of consolidated net profit (ETB 3.15bn). Per-share figure is approximate: total / 10,094,666 shares in issue at 30 Jun 2026. "
    f"Source: Wegagen notification to ESX dated 21 Sep 2026 ({ESX}/2026/09/33th-AGM-Resolution-Wegagen-Bank.pdf).", 2141446943.70 / 10094666)
add("WGBX", "agm", "33rd Ordinary General Assembly", "2026-09-19",
    "Approved FY2025/26 board performance and the audit report; declared the dividend; approved directors' fees; appointed TMS Plus Chartered Accountants as auditor for FY2026/27 (fee ETB 3,415,500 including VAT); elected an 11-member board for three years, subject to NBE approval. Audited statements were published in The Reporter on 19 Sep 2026.")
add("WGBX", "dividend", "Dividend declared on FY2024/25 profit (recorded in FY2025/26)", "2025-06-30",
    "ETB 1,944,300,000 declared (about 70% of FY2024/25 profit of ETB 2.78bn). Per share on the 7,028,357 shares in issue at 30 Jun 2025. Source: audited FY2025/26 consolidated statement of changes in equity.", 1944300000 / 7028357)
add("WGBX", "dividend", "Dividend declared on FY2023/24 profit (recorded in FY2024/25)", "2024-06-30",
    "ETB 1,246,781,000 declared. Per share on the 5,058,890 shares in issue at 30 Jun 2024. Source: audited FY2025/26 consolidated statement of changes in equity.", 1246781000 / 5058890)

# ---- Awash Bank (AWAB)
add("AWAB", "dividend", "Board-recommended dividend on FY2024/25 profit", "2025-06-30",
    "The Board recommended ETB 9,977,017,823 (about 53% of FY2024/25 net profit of ETB 18.71bn), payable in cash or as purchase of the bank's shares at each shareholder's choice. Per share is approximate: total / 27,874,601 fully paid shares. "
    "Source: Awash Bank Annual Report 2024/25 (30th AGM), section 2.5. A FY2025/26 dividend proposal had not been published when this was recorded.", 9977017823 / 27874601)
add("AWAB", "dividend", "Dividend declared on FY2023/24 profit (recorded in FY2024/25)", "2024-06-30",
    "ETB 6,159,439,000 declared (about 71% of FY2023/24 profit of ETB 8.68bn); ETB 1,008,056,000 paid in cash, with 5,041,459 new shares purchased through dividends. Per share on the 20,313,532 fully paid shares at 30 Jun 2024. Source: Awash Bank Annual Report 2024/25.", 6159439000 / 20313532)

# ---- Gadaa Bank (GDAB)
add("GDAB", "agm", "5th Ordinary General Meeting", "2026-09-19",
    "Held 19 Sep 2026 at the Addis International Convention Centre. Shareholders approved share transfers up to 4 Sep 2026, the board's and auditor's reports for FY2025/26, the recommendation for distribution of FY2025/26 net profit (amount not stated in the ESX disclosure), and directors' remuneration. "
    f"Source: Gadaa Bank disclosure to ESX dated 21 Sep 2026 ({ESX}/2026/09/Disclosure-of-5th-AGM-resolutions-Gadaa-Bank.pdf).")
add("GDAB", "dividend", "Dividend declared on FY2024/25 profit (recorded in FY2025/26)", "2025-06-30",
    "ETB 192.8m (about 56% of FY2024/25 profit of ETB 342.8m): ETB 165.2m was capitalised into share capital and ETB 8.3m into share premium (90% paid in new shares), ETB 16.5m payable in cash and ETB 2.9m dividend tax. Per share on the 1,303,918 shares at 30 Jun 2025. Source: audited FY2025/26 consolidated statement of changes in equity.", 192797000 / 1303918)
add("GDAB", "rights_issue", "Public offering and registration of additional shares in progress", "2026-06-29",
    "Gadaa Bank notified ESX that it is undertaking a public offering to raise capital and register additional existing shares, and was finalising regulatory prerequisites. Offer size and price not yet published. "
    f"Source: ESX notification dated 29 Jun 2026 ({ESX}/2026/07/Gadaa-bank-Additional-public-offer-notification.pdf).")

# ---- Abay Bank (ABAYB)
add("ABAYB", "dividend", "Cash dividends paid in FY2025/26", "2025-06-30",
    "Cash dividends paid during FY2025/26 were ETB 2,041,323,000 (FY2024/25: ETB 1,052,342,000), against net profit of ETB 3.00bn in FY2024/25. The declared total, including any part paid in shares, was not stated in the audited statements. Source: audited FY2025/26 cash flow statement.")

# ---- Sidama Bank (SIDAX)
add("SIDAX", "dividend", "Dividend on FY2024/25 profit: Board recommendation", "2025-06-30",
    "The Board recommended that dividends from FY2024/25 profit (appropriable net profit about ETB 72.3m) be paid in cash or used to buy new shares and complete earlier subscriptions, at each shareholder's choice; no amount stated. Cash dividends paid during FY2024/25 were ETB 45,686,433 (FY2023/24: ETB 35,021,700). Source: Sidama Bank Annual Report 2024/25.")

out = ["-- Generated by seeds/gen_dividends.py. Idempotent (replaces rows with the same company + title).", "BEGIN;"]
for t, typ, title, date, desc, amount in A:
    out.append(f"DELETE FROM corporate_actions WHERE title={q(title)} AND company_profile_id=(SELECT id FROM company_profiles WHERE ticker_symbol='{t}');")
    out.append(f"INSERT INTO corporate_actions (company_profile_id, action_type, title, event_date, description, amount, currency) SELECT id, '{typ}', {q(title)}, '{date}', {q(desc)}, {n(amount)}, 'ETB' FROM company_profiles WHERE ticker_symbol='{t}';")
out.append("COMMIT;")
open("004_dividends.sql", "w").write("\n".join(out) + "\n")
print(len(A), "corporate actions")
