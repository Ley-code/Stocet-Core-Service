#!/usr/bin/env python3
"""Generates seeds/002_esx_fundamentals.sql — real ESX listed-company data.

Sources (all public):
  - ESX company directory pages (profile, board, auditor, shares, filings): https://esx.et/directory/<slug>/
  - Audited financial statements for FY ended 30 June 2026 / 2025 (ESX PDFs)
  - Sidama Bank Annual Report 2024/25, Ethio Telecom prospectus Annex 1 (FY2022-24), Ethio Telecom FY2025/26 business report
Shares/capital = paid-up share capital at 30 Jun 2026 (audited) / par, which can exceed ESX "listed shares" after later capital increases.
Money values below are in statement units (Birr'000 unless noted) and converted to full ETB.
EPS is ETB per share at the stated par value (banks: 1,000; Ethio Telecom: 100).
"""
import json

def q(s):
    return "NULL" if s is None else "$q$" + s + "$q$"

def n(v):
    return "NULL" if v is None else repr(v)

ESX = "https://esx.et/wp-content/uploads"

companies = [
 dict(t="WGBX", name="Wegagen Bank Share Company", sector="Banking", nob="Financial Services - Banking", listed="2025-01-10", inc="1997-06-11",
      auditor="Tafesse, Shisema, and Ayelew Certified Audit Partnership (FY2025); TMS Plus Audit Service LLP (FY2026)", shares=10094666, capital=10094666000, par=1000,
      email="info@wegagenbanksc.com.et", phone="+251-115-523800", addr="Wegagen Tower, Ras Mekonnen Street, Kirkos Sub-City, Woreda 10, Addis Ababa (in front of Addis Ababa Stadium)", web="https://www.wegagen.com", holders=14549),
 dict(t="GDAB", name="Gadaa Bank Share Company", sector="Banking", nob="Financial Services - Banking", listed="2025-06-13", inc="2022-04-26",
      auditor="Tewodros and Fikre Audit Services Partnership (FY2025); Kokeb & Melkamu Audit Service LLP (FY2026)", shares=1654458, capital=1654458000, par=1000,
      email="info@gadaabank.com.et", phone="+251-116-392578", addr="Head Office, Gotera, Kirkos Sub-City, Woreda 03, House No. 745, Addis Ababa", web="https://gadaabank.com.et", holders=31136),
 dict(t="AWAB", name="Awash Bank Share Company", sector="Banking", nob="Financial Services - Banking", listed="2026-04-23", inc="1994-11-10",
      auditor="Audit Service Corporation (ASC)", shares=38492985, capital=38492985000, par=1000,
      email="contactcenter@awashbank.com", phone="0115-57-12-54 / 0115-57-13-24", addr="Ras Abebe Aregay Street, P.O. Box 12638, Addis Ababa", web="https://awashbank.com", holders=12143),
 dict(t="TELE", name="Ethio Telecom Share Company", sector="Telecom", nob="Telecoms, Infrastructure Sharing, and Mobile Financial Services", listed="2026-05-26", inc="2024-07-01",
      auditor="Audit Services Corporation", shares=1000000000, capital=100000000000, par=100,
      email="994@ethionet.et", phone="+251-115-510-500", addr="Lideta Sub-City, Woreda 10, P.O. Box 1047, Addis Ababa", web="https://www.ethiotelecom.et", holders=47305),
 dict(t="ABAYB", name="Abay Bank Share Company", sector="Banking", nob="Financial Services - Banking", listed="2026-06-25", inc="2010-07-14",
      auditor="Tewodros and Fikre Audit Services Partnership", shares=9657286, capital=9657286000, par=1000,
      email="info@abaybanksc.com", phone="+251-11554-9735", addr="Ras Abebe Aregay Street, P.O. Box 12638, Addis Ababa", web="https://www.abaybanksc.com", holders=4738),
 dict(t="BOAX", name="Bank of Abyssinia Share Company", sector="Banking", nob="Financial Services - Banking", listed="2026-07-28", inc="1996-02-16",
      auditor="Tewodros and Fikre Audit Services Partnership", shares=17481120, capital=17481120000, par=1000,
      email="ContactCenterSupport@bankofabyssinia.com", phone="+251-11551-4130", addr="Kirkos Sub-City, Woreda 07, House No. 351/01, P.O. Box 12947, Addis Ababa", web="https://www.bankofabyssinia.com", holders=4974),
 dict(t="SIDAX", name="Sidama Bank Share Company", sector="Banking", nob="Financial Services - Banking", listed="2026-09-28", inc="2022-07-01",
      auditor="Tewodros and Fikre Audit Services Partnership", shares=1238763, capital=1238763000, par=1000,
      email="info@sidamabanksc.com", phone="+251 114704726 / +251 114706740", addr="Kirkos Sub-city, Woreda 04, House No. 892, Sierra Leone Avenue, Lancha Mekor Plaza, 3rd Floor, P.O. Box 1157, Addis Ababa", web="https://www.sidamabanksc.com", holders=3257),
]

board = {
 "AWAB": [("Gure Kumssa Deressa","Chairman"),("Kedeme Teshome Woldeyes","Vice Chairman"),("Abaynesh Merga Jena","Board Director"),("Arega Gella Hedu","Board Director"),("Beza Tadesse Gurmu","Board Director"),("Elias Geneti Sima","Board Director"),("Frehiwot Alemayehu Kelbessa","Board Director"),("Hordofa Bekele Gumie (Dr. of Honor)","Board Director"),("Million Duguma Hundie","Board Director"),("Roba Arega Kitata","Board Director"),("Yoseph Yonas Kumsa","Board Director"),("Tsehay Shiferaw Ensermu","Chief Executive Officer")],
 "SIDAX": [("Abrham Marshalo Etiso","Board Chairperson"),("Tsegaye Anebo Alemaso","Board Vice Chairperson"),("Timotwos Hayeso Borsamo","Board Member"),("Dr. Wogene Markos Dumo","Board Member"),("Woinshet Mengesha Sinkale","Board Member"),("Muluwork Derese Degisew","Board Member"),("Haile Harsisa Haro","Board Member"),("Shimelis Hatiso Hankalo","Board Member"),("Shebeku Magane Gafa","Board Member")],
 "BOAX": [("Mekonnen Manyazewal Meka","Chairperson of the Board"),("Aemero Belete Simegne","Deputy Chairperson"),("Emebet Woldeher Yezengaw","Director"),("Kassahun Zewdie Mengesha","Director"),("Meseret Melese Tefera","Director"),("Molalign Melese Mengistu","Director"),("Solomon Alula Awlachew","Director"),("Yerom Gessesse Yeneneh","Director"),("Dr. Yihenew Zewdie Lemma","Director"),("Bekalu Zeleke Ewnetu","Chief Executive Officer"),("Abreham Gebeyehu G/Hanna","Chief Financial Officer")],
 "ABAYB": [("Dr. Amlaku Asres Zewdie","Board Chairperson"),("Tadesse Assefa Tiruneh","Board Vice Chairperson"),("Abrham Ayalew Ejigu","Board Member"),("Agerie Tigabu Mekonnen","Board Member"),("Demis Teshager Abeje","Board Member"),("Jemberie Workneh Temesgen","Board Member"),("Mulat Tsega Ayele","Board Member"),("Tibebu Alula Getahun","Board Member"),("Martha Ayenew Dagnachew","Board Member"),("Yehuala Gessesse Gelaye","Chief Executive Officer"),("Desalegn Ayalew Worku","Chief Corporate Service Officer")],
 "TELE": [("H.E Temesgen Tiruneh","Board Chairman"),("H.E Dr. Eyob Tekalign","Board Deputy Chairman"),("H.E Ato Worku Gachena","Board Member"),("H.E Ato Muluneh Desalegn","Board Member"),("H.E Ato Yodahe Arayaselassie","Board Member"),("Ato Kasahun Seboka","Board Member")],
 "GDAB": [("Dr. Hassan Hussien Kedir","Board Chairperson"),("Mr. Hailu Ifa Gonfa","Board Vice Chairperson"),("Dr. Degefa Duresa Obbo","Board Member"),("Eng. Abdo Galeto Anota","Board Member"),("Dr. Birhanu Assefa Demissie","Board Member"),("Mr. Shiferaw Rufie Bodo","Board Member"),("Mr. Wassihun Amenu Tiyiti","Board Member"),("Mr. Hamdeno Mideso Woya","Board Member"),("Mrs. Semira Mohammed Abdella","Board Member"),("Mr. Alazar Adula Yatene","Board Member"),("Dr. Gutu Teso Boka","Board Member"),("Wolde Bulto Adugna","Chief Executive Officer")],
 "WGBX": [("Ato Abdishu Hussien","Board Chairperson"),("Ato Woldegabriel Naizghi","Board Vice Chairperson"),("Ato Fikru Jiregna","Board Member"),("Ato Fithanegest Gebru","Board Member"),("Ato Gebreegziabher Hadush","Board Member"),("Ato Hassen Yesuf","Board Member"),("Ato Surfel Berhe Weldu","Board Member"),("Ato Tesfatsion Desta Tesfay","Board Member"),("Ato Alemseged Assefa Abera","Board Member"),("Ato Zenfu Asfaw Gebretinsae","Board Member"),("Dr. Aklilu Wubet","Chief Executive Officer")],
}

# unit multiplier: statement unit -> ETB
K = 1000
def fr(t, period, unit, **kw):
    kw["t"] = t; kw["period"] = period; kw["unit"] = unit
    return kw

reports = [
 # Awash — separate (bank-only) statements, Birr'000
 fr("AWAB","2026-06-30",K,toi=69172513,pbt=38312238,np=29341603,ta=624106921,tl=533713795,te=90393125,dep=469495118,loans=268477527,eps=858.0,ext=dict(basis="separate",net_interest_income=38828098*K,fees_commission_income=22910491*K,total_income_note="Total operating income",source=f"{ESX}/2026/10/Awash-Bank-Audited-Financial-Statement-June-30-2026.pdf")),
 fr("AWAB","2025-06-30",K,toi=54769402,pbt=25668278,np=18712623,ta=442596547,tl=382612041,te=59984506,dep=332903235,loans=215030059,eps=783.0,ext=dict(basis="separate",net_interest_income=26287891*K,fees_commission_income=17241066*K,source=f"{ESX}/2026/10/Awash-Bank-Audited-Financial-Statement-June-30-2026.pdf (comparative)")),
 # Abay
 fr("ABAYB","2026-06-30",K,toi=13120143,pbt=6222965,np=4545938,ta=119687868,tl=102110260,te=17577608,dep=92739547,loans=61152836,eps=524.54,ext=dict(basis="separate",net_interest_income=8614360*K,fees_commission_income=2709456*K,source=f"{ESX}/2026/10/Abay-Bank-Audited-Financial-Statement-June-30-2026.pdf")),
 fr("ABAYB","2025-06-30",K,toi=10476832,pbt=4215729,np=2999863,ta=91345765,tl=79136414,te=12209351,dep=71914463,loans=48648696,eps=463.31,ext=dict(basis="separate",net_interest_income=6263828*K,fees_commission_income=2038082*K,source=f"{ESX}/2026/10/Abay-Bank-Audited-Financial-Statement-June-30-2026.pdf (comparative)")),
 # Bank of Abyssinia (deposits incl. interest-free deposits)
 fr("BOAX","2026-06-30",K,toi=40291549,pbt=15849406,np=11535777,ta=381497212,tl=341985031,te=39512182,dep=277982151+40803442,loans=248997369,eps=770.0,ext=dict(basis="separate",net_interest_income=28743103*K,interest_income=43689742*K,branches=981,source=f"{ESX}/2026/10/Annual-Report-2025-26-BOA.pdf",note="Deposits include interest-free deposits (Birr 40.8bn). EPS reported as 77% of Birr 1,000 par.")),
 fr("BOAX","2025-06-30",K,toi=28748312,pbt=10101095,np=7298505,ta=286226976,tl=257385045,te=28841931,dep=211654830+31524908,loans=193379151,eps=500.0,ext=dict(basis="separate",net_interest_income=22395569*K,interest_income=32718727*K,source=f"{ESX}/2026/10/Annual-Report-2025-26-BOA.pdf (comparative)")),
 # Gadaa — consolidated
 fr("GDAB","2026-06-30",K,toi=2309057,pbt=1045989,np=759231,ta=16793877,tl=14019054,te=2774823,dep=11610198,loans=5842957,eps=520.0,ext=dict(basis="consolidated",net_interest_income=653622*K,fees_commission_income=1171743*K,source=f"{ESX}/2026/10/Gadaa-Bank-Audited-Financial-Statement-June-30-2026.pdf",note="Reported EPS Birr 0.52 per Birr 1 of capital = 520 per Birr 1,000 share.")),
 fr("GDAB","2025-06-30",K,toi=1415538,pbt=444790,np=342756,ta=10139499,tl=8359989,te=1779511,dep=7082923,loans=3658246,eps=290.0,ext=dict(basis="consolidated",net_interest_income=407969*K,fees_commission_income=637588*K,source=f"{ESX}/2026/10/Gadaa-Bank-Audited-Financial-Statement-June-30-2026.pdf (comparative)",note="Source statement has a Birr 100m footing difference between assets and liabilities+equity; total assets taken from liabilities+equity.")),
 # Wegagen — consolidated. TOI = net interest + net fee/commission + other operating income (before impairment)
 fr("WGBX","2026-06-30",K,toi=11626020,pbt=4306596,np=3154253,ta=111316522,tl=93971006,te=17345516,dep=73517597,loans=64590783,eps=344.0,ext=dict(basis="consolidated",net_interest_income=7063059*K,fees_commission_income=3477954*K,source=f"{ESX}/2026/10/Wegagen-Bank-Audited-Financial-Statement-June-30-2026.pdf",note="Reported EPS 34.40% of Birr 1,000 par.")),
 fr("WGBX","2025-06-30",K,toi=10223931,pbt=3852553,np=2777510,ta=84337595,tl=71529622,te=12807973,dep=61651819,loans=51260502,eps=461.0,ext=dict(basis="consolidated",net_interest_income=5303178*K,fees_commission_income=2464927*K,source=f"{ESX}/2026/10/Wegagen-Bank-Audited-Financial-Statement-June-30-2026.pdf (comparative)")),
 # Sidama — statements in whole Birr
 fr("SIDAX","2025-06-30",1,toi=500295914,pbt=90787638,np=72314055,ta=5519426992,tl=4208328529,te=1311098463,dep=3674915745,loans=2390120571,eps=75.0,ext=dict(basis="separate",net_interest_income=248558913,fees_commission_income=81928100,source=f"{ESX}/2026/10/Sidama-Annual-Report-2024-2025.pdf",note="Reported EPS 7.5% of Birr 1,000 par. FY2025/26 audited statements not yet published on ESX.")),
 fr("SIDAX","2024-06-30",1,toi=308327018,pbt=29426888,np=27135894,ta=2620268280,tl=1531455616,te=1088812663,dep=1224875020,loans=1112924098,eps=33.0,ext=dict(basis="separate",net_interest_income=234376166,fees_commission_income=26104866,source=f"{ESX}/2026/10/Sidama-Annual-Report-2024-2025.pdf (comparative)")),
 # Ethio Telecom
 fr("TELE","2026-06-30",K,toi=215820000,pbt=None,np=None,ta=None,tl=None,te=None,dep=None,loans=None,eps=None,ext=dict(basis="management report",revenue_growth_pct=33.2,total_customers=90120000,active_mobile_subscribers=72090000,active_mobile_data_users=51530000,telebirr_customers=60600000,telebirr_transaction_value_etb=4.19e12,telebirr_transactions=2.61e9,fixed_broadband_customers=1000000,smartphone_penetration_pct=45.2,source="https://www.stockmarket.et/ethio-telecom-revenue-climbs-33-2-to-etb-215-8-billion-but-annual-report-reveals-key-growth-targets-missed/",note="Revenue and operating KPIs from the FY2025/26 business performance report. Profit and balance sheet will be added when audited FY2025/26 statements are filed on ESX.")),
 fr("TELE","2024-06-30",K,toi=91371290,pbt=29790459,np=19013443,ta=214191635,tl=100698117,te=113493516,dep=None,loans=None,eps=19.01,ext=dict(basis="separate",ebitda=41220047*K,operating_profit=32858411*K,source=f"{ESX}/2025/06/Ethio-Telecom-Share-Company-Three-years-historical-financial-information-2022-2024.pdf",note="EPS = profit after tax / 1,000,000,000 shares (Birr 100bn capital at Birr 100 par).")),
 fr("TELE","2023-06-30",K,toi=71525993,pbt=28186778,np=18028014,ta=199659771,tl=95379674,te=104280097,dep=None,loans=None,eps=18.03,ext=dict(basis="separate (restated)",ebitda=37619767*K,operating_profit=31122907*K,source=f"{ESX}/2025/06/Ethio-Telecom-Share-Company-Three-years-historical-financial-information-2022-2024.pdf (comparative)")),
]

def full(v, unit):
    return None if v is None else v * unit

filings = {
 "AWAB": [("Prospectus","Awash Bank S.C. Prospectus","2026-04-01",f"{ESX}/2026/04/Awash-Bank-SC-Prospectus.pdf"),
          ("Annual Report","Audited Financial Statements — year ended 30 June 2026","2026-10-01",f"{ESX}/2026/10/Awash-Bank-Audited-Financial-Statement-June-30-2026.pdf"),
          ("Material Disclosure","Provisional annual key financial highlights FY2025/26","2026-07-01",f"{ESX}/2026/07/Awash-Bank-Financial-Performance-Announcement.pdf"),
          ("Material Disclosure","Closed period announcement","2026-07-01",f"{ESX}/2026/07/Closed-Period-Announcement-Awash-Bank.pdf"),
          ("Material Disclosure","Appointment of new Board of Directors","2026-07-01",f"{ESX}/2026/07/Awash-Bank-Appointment-of-New-Board-of-Directors.pdf")],
 "ABAYB": [("Prospectus","Abay Bank S.C. Prospectus","2026-06-01",f"{ESX}/2026/06/Prospectus_Abay_Bank.pdf"),
           ("Annual Report","Audited Financial Statements — year ended 30 June 2026","2026-10-01",f"{ESX}/2026/10/Abay-Bank-Audited-Financial-Statement-June-30-2026.pdf"),
           ("Annual Report","Audited Financial Statements FY2025","2026-06-01",f"{ESX}/2026/06/Audited-Financial-Statement-FY-2025.pdf"),
           ("Annual Report","Audited Financial Statements FY2024","2026-06-01",f"{ESX}/2026/06/Audited-Financial-Statement-FY-2024.pdf"),
           ("Interim Report","Auditor review interim report, July–December 2025","2026-06-01",f"{ESX}/2026/06/Auditor_Review_Interim_Report_July_2025_December_2025.pdf"),
           ("Material Disclosure","Closed period announcement","2026-07-01",f"{ESX}/2026/07/Closed-Period-Announcement-Abay-Bank.pdf"),
           ("Material Disclosure","Board chairman and assistant chairman appointment","2026-07-01",f"{ESX}/2026/07/Abay-Bank_Board-chairman-Assistant-boardchairman.pdf")],
 "BOAX": [("Prospectus","Bank of Abyssinia S.C. Prospectus","2026-07-01",f"{ESX}/2026/07/Bank-of-Abyssinia-Share-Company-Prospectus-.pdf"),
          ("Annual Report","Annual Report 2025/26","2026-10-01",f"{ESX}/2026/10/Annual-Report-2025-26-BOA.pdf"),
          ("Annual Report","Audited Financial Statements — year ended 30 June 2026","2026-10-01",f"{ESX}/2026/10/BOA-Audited-Financial-Statement-June-30-2026.pdf"),
          ("Annual Report","Annual Report 2024/25","2026-07-01",f"{ESX}/2026/07/2024-25-BOA-Annual-report-final-print-ready-2-3.pdf"),
          ("Annual Report","Annual Report 2023/24","2026-07-01",f"{ESX}/2026/07/BOA-2023-24-ANNUAL-REPORT-2.pdf"),
          ("Interim Report","Interim Financial Statements — 31 December 2025","2026-07-01",f"{ESX}/2026/07/Interim-Financial-Statments-December-31-2025.pdf"),
          ("AGM Notice","AGM Notification 2026","2026-10-01",f"{ESX}/2026/10/BOA-AGM-Notification-2026.pdf"),
          ("Material Disclosure","Notification of 17th AGM resolutions","2026-10-01",f"{ESX}/2026/10/Notification-of-17th-AGM-Resolution-BOA.pdf")],
 "SIDAX": [("Prospectus","Sidama Bank S.C. approved prospectus (2026)","2026-10-01",f"{ESX}/2026/10/Sidama-Approved-Prospectus.pdf"),
           ("Interim Report","Reviewed interim financial statements","2026-10-01",f"{ESX}/2026/10/Sidama-Bank-Reviewed-Interim.pdf"),
           ("Annual Report","Annual Report 2024/25","2026-10-01",f"{ESX}/2026/10/Sidama-Annual-Report-2024-2025.pdf"),
           ("Annual Report","Annual Report 2023/24","2026-10-01",f"{ESX}/2026/10/Sidama-Annual-Report-2023-2024.pdf"),
           ("Annual Report","Annual Report 2022/23","2026-10-01",f"{ESX}/2026/10/Sidama-Annual_report_2022-23.pdf")],
 "TELE": [("Prospectus","Ethio Telecom S.C. Prospectus","2025-06-01",f"{ESX}/2025/06/Ethio-Telecom-Share-Company-Prospectus.pdf"),
          ("Annual Report","Historical financial information FY2022–FY2024 (audited)","2025-06-01",f"{ESX}/2025/06/Ethio-Telecom-Share-Company-Three-years-historical-financial-information-2022-2024.pdf"),
          ("Material Disclosure","Notification of closed period (30/6/2026)","2026-07-01",f"{ESX}/2026/07/Notfication-of-clouse-priod-ESX-30_6_2026-Ethio-Telecom.pdf"),
          ("Material Disclosure","Resignation of CEO","2026-10-01",f"{ESX}/2026/10/Ethio-Tele-Resignation-of-CEO.pdf"),
          ("Material Disclosure","Appointment of Chairperson","2026-10-01",f"{ESX}/2026/10/Ethio-Tele-Appointment-of-Chairperson.pdf")],
 "GDAB": [("Prospectus","Gadaa Bank final prospectus dated 17 June 2025","2025-06-01",f"{ESX}/2025/06/GB-Final-Propsectus-Dated-June-17-2025.pdf"),
          ("Annual Report","Audited Financial Statements — year ended 30 June 2026","2026-10-01",f"{ESX}/2026/10/Gadaa-Bank-Audited-Financial-Statement-June-30-2026.pdf"),
          ("Annual Report","Audited Financial Statements — year ended 30 June 2025","2025-11-01",f"{ESX}/2025/11/Audited-Financial-Statement-for-the-year-ended-30-June-2025.pdf"),
          ("Annual Report","Annual Report 2023/24","2025-06-01",f"{ESX}/2025/06/Gadaa-Bank-2023-24-Annual-Report.pdf"),
          ("Annual Report","Annual Report 2022/23","2025-06-01",f"{ESX}/2025/06/Gadaa-Bank-Annual-Report-Final-2022-23-_compressed.pdf"),
          ("Material Disclosure","Additional public offer notification","2026-07-01",f"{ESX}/2026/07/Gadaa-bank-Additional-public-offer-notification.pdf"),
          ("AGM Notice","5th AGM notification","2026-08-01",f"{ESX}/2026/08/Gadaa-Bank-AGM-Notification.pdf"),
          ("Material Disclosure","Disclosure of 5th AGM resolutions","2026-09-01",f"{ESX}/2026/09/Disclosure-of-5th-AGM-resolutions-Gadaa-Bank.pdf")],
 "WGBX": [("Prospectus","Wegagen Bank prospectus for registration of shares currently held by shareholders","2025-03-01",f"{ESX}/2025/03/Wegagen_Bank_S_C_Prospectus-for-Registration-of-Shares-Currently-Held-by-Shareholders.pdf"),
          ("Annual Report","Audited Financial Statements — year ended 30 June 2026","2026-10-01",f"{ESX}/2026/10/Wegagen-Bank-Audited-Financial-Statement-June-30-2026.pdf"),
          ("Annual Report","Annual Report 2024/25","2025-10-01",f"{ESX}/2025/10/Wegagen_Bank_2024_25_Annual_Report_4c4ff4570b.pdf"),
          ("Annual Report","Annual Report 2023/24","2025-03-01",f"{ESX}/2025/03/Wegagen_Bank_Annual_Report_2023-24.pdf"),
          ("Material Disclosure","Financial highlights announcement FY2025/26","2026-07-01",f"{ESX}/2026/07/Wegagen-Bank-Announcment-Finacial-Highlight.pdf"),
          ("AGM Notice","Announcement of Annual General Meeting of Shareholders","2026-08-01",f"{ESX}/2026/08/Wegagen-Bank-Announcement-of-Annual-General-Meeting-of-Shareholders.pdf"),
          ("Material Disclosure","33rd AGM resolutions","2026-09-01",f"{ESX}/2026/09/33th-AGM-Resolution-Wegagen-Bank.pdf")],
}

out = ["-- Generated by seeds/gen_fundamentals.py — real ESX data (see script header for sources). Idempotent.", "BEGIN;"]
for c in companies:
    out.append(f"""INSERT INTO company_profiles (ticker_symbol, company_name, sector, nature_of_business, security_type, listing_date, incorporation_date, auditor, shares, capital, contact_email, contact_phone, contact_address, website)
VALUES ({q(c['t'])},{q(c['name'])},{q(c['sector'])},{q(c['nob'])},'Equity','{c['listed']}','{c['inc']}',{q(c['auditor'])},{c['shares']},{c['capital']},{q(c['email'])},{q(c['phone'])},{q(c['addr'])},{q(c['web'])})
ON CONFLICT (ticker_symbol) DO UPDATE SET company_name=EXCLUDED.company_name, sector=EXCLUDED.sector, nature_of_business=EXCLUDED.nature_of_business, security_type=EXCLUDED.security_type,
 listing_date=EXCLUDED.listing_date, incorporation_date=EXCLUDED.incorporation_date, auditor=EXCLUDED.auditor, shares=EXCLUDED.shares, capital=EXCLUDED.capital,
 contact_email=EXCLUDED.contact_email, contact_phone=EXCLUDED.contact_phone, contact_address=EXCLUDED.contact_address, website=EXCLUDED.website, updated_at=now();""")

for t, members in board.items():
    out.append(f"DELETE FROM board_members WHERE company_profile_id=(SELECT id FROM company_profiles WHERE ticker_symbol='{t}');")
    vals = ",".join(f"({q(nm)},{q(role)})" for nm, role in members)
    out.append(f"INSERT INTO board_members (company_profile_id, name, role) SELECT p.id, v.name, v.role FROM company_profiles p, (VALUES {vals}) AS v(name, role) WHERE p.ticker_symbol='{t}';")

for r in reports:
    u = r["unit"]
    ext = dict(r["ext"]); ext["currency"] = "ETB"; ext["fiscal_year_end"] = r["period"]
    par = next(c["par"] for c in companies if c["t"] == r["t"])
    ext["par_value"] = par
    ext["eps_basis"] = f"ETB per share (par {par})"
    out.append(f"""INSERT INTO financial_reports (company_profile_id, period_end_date, net_profit, total_operating_income, profit_before_tax, total_assets, total_equity, total_liabilities, eps, customer_deposits, loans_and_advances, extended_metrics)
SELECT p.id,'{r['period']}',{n(full(r['np'],u))},{n(full(r['toi'],u))},{n(full(r['pbt'],u))},{n(full(r['ta'],u))},{n(full(r['te'],u))},{n(full(r['tl'],u))},{n(r['eps'])},{n(full(r['dep'],u))},{n(full(r['loans'],u))},{q(json.dumps(ext))}::jsonb
FROM company_profiles p WHERE p.ticker_symbol='{r['t']}'
ON CONFLICT (company_profile_id, period_end_date) DO UPDATE SET net_profit=EXCLUDED.net_profit, total_operating_income=EXCLUDED.total_operating_income, profit_before_tax=EXCLUDED.profit_before_tax,
 total_assets=EXCLUDED.total_assets, total_equity=EXCLUDED.total_equity, total_liabilities=EXCLUDED.total_liabilities, eps=EXCLUDED.eps, customer_deposits=EXCLUDED.customer_deposits,
 loans_and_advances=EXCLUDED.loans_and_advances, extended_metrics=EXCLUDED.extended_metrics, updated_at=now();""")

# older Wegagen / Gadaa rows: only correct EPS scale (was % / fraction)
out.append("UPDATE financial_reports SET eps=368.9, extended_metrics = extended_metrics || '{\"eps_basis\":\"ETB per share (par 1000)\",\"par_value\":1000}'::jsonb WHERE company_profile_id=(SELECT id FROM company_profiles WHERE ticker_symbol='WGBX') AND period_end_date='2024-06-30';")
out.append("UPDATE financial_reports SET eps=80.0, extended_metrics = extended_metrics || '{\"eps_basis\":\"ETB per share (par 1000)\",\"par_value\":1000}'::jsonb WHERE company_profile_id=(SELECT id FROM company_profiles WHERE ticker_symbol='GDAB') AND period_end_date='2024-06-30';")

for t, rows in filings.items():
    out.append(f"DELETE FROM company_filings WHERE company_profile_id=(SELECT id FROM company_profiles WHERE ticker_symbol='{t}') AND document_url LIKE 'https://esx.et/%';")
    vals = ",".join(f"({q(ft)},{q(title)},{q(url)},'{d}T00:00:00Z'::timestamptz)" for ft, title, d, url in rows)
    out.append(f"INSERT INTO company_filings (company_profile_id, filing_type, title, description, document_url, filed_at) SELECT p.id, v.ft, v.title, 'Published on the ESX company directory.', v.url, v.d FROM company_profiles p, (VALUES {vals}) AS v(ft,title,url,d) WHERE p.ticker_symbol='{t}';")

# real listing events (replaces earlier placeholder AWAB 'General Listiing IPO')
out.append("DELETE FROM corporate_actions WHERE action_type='listing' AND company_profile_id IN (SELECT id FROM company_profiles);")
for c in companies:
    desc = f"{c['name']} began trading on the ESX Main Market under ticker {c['t']}."
    out.append(f"INSERT INTO corporate_actions (company_profile_id, action_type, title, event_date, description) SELECT id,'listing',{q('Listed on ESX Main Market ('+c['t']+')')},'{c['listed']}',{q(desc)} FROM company_profiles WHERE ticker_symbol='{c['t']}';")
out.append("COMMIT;")
open("002_esx_fundamentals.sql","w").write("\n".join(out)+"\n")
print("wrote", len(out), "statements")
