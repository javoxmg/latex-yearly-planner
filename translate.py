"""Translate the generated out/*.tex files into another language.

Usage (normally run by single.sh when TRANSLATION is set):

    python3 translate.py spanish

It loads translations/<language>.json, a flat {"English": "Translation"}
dictionary, and replaces every English label in the generated pages with
its translation. Two rules keep the hyperlinks working:

* Link *targets* are never translated. The first argument of every
  \\hyperlink{...}{...} and \\hypertarget{...}{...} is masked before the
  replacement pass and restored afterwards, so a page that links to
  "September" keeps pointing at the "September" anchor while its visible
  text becomes "Septiembre".
* Everything else is replaced as whole words (case sensitive), longest
  key first, so "Notes Index" wins over "Notes", and macro names such as
  \\myNumWeeklyLines are left alone (a backslash or a letter right before
  the word disables the match).

Keys missing from the JSON simply stay in English.
"""

import json
import re
from glob import glob
from sys import argv

language = argv[1].lower()
TRANSLATION_FOLDER = "translations/"
file = f"{TRANSLATION_FOLDER}{language}.json"

if file in glob(f"{TRANSLATION_FOLDER}*.json"):
    with open(file, "r", encoding="utf-8") as f:
        translation = json.load(f)
else:
    raise ValueError("Requested translation is not currently supported.\nThe program will now exit.")

print(f"Translating pdf to {language}")

FILES = [
    "out/annual.tex",
    "out/quarterly.tex",
    "out/monthly.tex",
    "out/weekly.tex",
    "out/daily.tex",
    "out/daily_reflect.tex",
    "out/daily_notes.tex",
    "out/notes_indexed.tex",
]

# Keys that contain characters other than letters/spaces/apostrophes (for
# example the weekday-letter row "W & M & T & W & T & F & S & S") are
# replaced literally; the rest as whole words.
literal_keys = [k for k in translation if not re.fullmatch(r"[A-Za-z' ]+", k) and k != "May (short)"]
word_keys = sorted((k for k in translation if k not in literal_keys), key=len, reverse=True)
word_re = re.compile(r"(?<![A-Za-z\\])(" + "|".join(re.escape(k) for k in word_keys) + r")(?![A-Za-z])")

LINK_RE = re.compile(r"(\\(?:hyperlink|hypertarget)\{)([^{}]*)(\})")

# The months-on-side layout prints the month tabs with the short month
# name ("Sep", "Oct"...). For May the short and the long name are the same
# word, so inside the tab table the display "May" is translated with the
# "May (short)" key (falling back to the plain "May" key) instead.
SIDE_TABS_RE = re.compile(
    r"\\begin\{tabularx\}\{\\myLenHeaderSideMonthsWidth\}.*?\\end\{tabularx\}", re.DOTALL
)
MAY_SHORT = "\x02MAYSHORT\x02"


def translate_text(text: str) -> str:
    targets: list[str] = []

    def mask(m: re.Match) -> str:
        targets.append(m.group(2))
        return f"{m.group(1)}\x00{len(targets) - 1}\x00{m.group(3)}"

    text = SIDE_TABS_RE.sub(lambda m: m.group(0).replace(r"\hyperlink{May}{May}", r"\hyperlink{May}{" + MAY_SHORT + "}"), text)
    text = LINK_RE.sub(mask, text)

    for key in literal_keys:
        text = text.replace(key, translation[key])

    text = word_re.sub(lambda m: translation[m.group(1)], text)
    text = text.replace(MAY_SHORT, translation.get("May (short)", translation.get("May", "May")))

    return re.sub(r"\x00(\d+)\x00", lambda m: targets[int(m.group(1))], text)


for path in FILES:
    try:
        with open(path, "r", encoding="utf-8") as fh:
            original = fh.read()
    except FileNotFoundError:
        continue

    with open(path, "w", encoding="utf-8") as fh:
        fh.write(translate_text(original))
