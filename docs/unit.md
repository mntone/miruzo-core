# manbytes — A culturally-inspired yet practical size unit

## 1. Background: Japan’s 10<sup>4</sup>-based grouping

In Japanese (and broadly in East Asian numeral systems), numbers are commonly
grouped in **units of four digits (10<sup>4</sup>)** rather than three. This
system is deeply embedded in everyday language, using distinct magnitude names
such as:

- man (万, 10<sup>4</sup>)
- oku (億, 10<sup>8</sup>)
- chō (兆, 10<sup>12</sup>)
- kei (京, 10<sup>16</sup>)

Expressions like "1 man", "12 man", or "300 man" are immediately understood,
and changes in magnitude are perceived intuitively as relative weight rather
than as raw digits.

By contrast, Western numbering systems are based on **three-digit (10<sup>3</sup>)
groupings**, with units such as:

- thousand (10<sup>3</sup>)
- million (10<sup>6</sup>)
- billion (10<sup>9</sup>)

While this difference is often viewed as cultural, it also has practical
implications when numbers are presented in user interfaces.

In the context of image processing and UI design, using 10<sup>4</sup> bytes
(approximately 10 KB) as a base unit turns out to be a particularly effective
scale:

- it is not too fine-grained,
- yet not overly coarse,

making it well-suited for reasoning about perceived image size and loading cost.


## 2. Why not kilobytes or megabytes?

Conventional byte-based units such as kilobytes and megabytes are standard from
an implementation perspective, but they are often ill-suited as inputs for
UX-level heuristics.

**kilobytes (10³ bytes)**

Kilobytes are *too fine-grained*. Small numerical differences introduce noise
when the goal is to classify image size or loading cost roughly.

For example, while the difference between 120 KB and 180 KB is numerically
clear, in terms of rendering strategy or perceived load time, these sizes
usually belong to the same category.

**megabytes (10<sup>6</sup> bytes)**

Megabytes, on the other hand, are *too coarse*. Differences such as 1.1 MB
versus 1.9 MB can result in noticeably different user experiences on mobile
networks, yet they are often flattened into the same "~2 MB" bucket.

**kibibytes, mebibytes, and similar units (2<sup>10</sup>, 2<sup>20</sup> bytes)**

Units such as kibibytes (KiB) and mebibytes (MiB) are technically precise.
However, that level of precision is unnecessary for UX-oriented decisions such
as progress indicators, preload thresholds, or branching logic.

In these contexts, **stable categorization matters more than exactness.**


## 3. Why 10<sup>4</sup> works so well for images

When deciding image loading and rendering strategies, the exact byte count
matters less than the *perceived weight* of the image.

In most UIs, image size influences decisions such as:

- whether it can be shown immediately,
- whether it feels fast to load,
- whether a progress indicator is needed,
- or whether special handling is required.

In practice, these decisions tend to fall into ranges like:

- 0–100 KB: immediate display
- 200–400 KB: fast loading
- 1–2 MB: ~0.3–1.0 s on mobile networks
- 4–8 MB: progress UI recommended
- 10 MB and above: special handling

Using a 10<sup>4</sup>-byte unit (approximately 10 KB) allows these ranges to
be expressed in a smooth and continuous scale.

For example:

- 28 manbytes → ~280 KB
- 120 manbytes → ~1.2 MB
- 480 manbytes → ~4.8 MB

In this form, the numeric value itself becomes a scale that can be used directly
for UX-level decisions.


## 4. The manbytes unit

Manbytes is a size unit defined as 10<sup>4</sup> bytes.

**1 manbyte = 10<sup>4</sup> bytes**

In miruzo, raw file sizes are normalized into this unit. The conversion is
defined as follows:

**manbytes = ceil(bytes / 10_000)**

This ensures that even small files are assigned a meaningful, non-zero size
class, which stabilizes UX-level decisions.

Manbytes is **not intended for precise storage reporting**. Its primary use
cases include:

- branching image loading strategies,
- deciding whether to show a progress indicator,
- preload and lazy-load heuristics,
- determining rendering priority.

Whenever exact byte counts are required, traditional units such as bytes,
kilobytes, or megabytes remain available.

### Implementation notes

Code that uses manbytes assumes the definitions and conversion formula
specified in this section.

See the following implementations:

- [`miruzo-core/app/models/api/images/variant.py`](../app/models/api/images/variant.py)
- [`miruzo-core/app/models/api/utils/units.py`](../app/models/api/utils/units.py)

Using `ceil` ensures that a value of `0` unambiguously represents a zero-byte
file, rather than a non-zero file rounded down to zero.


## 5. Future extensions: okubytes & chobytes

Beyond 10<sup>4</sup> (man), the Japanese numeric system defines larger units
such as 10<sup>8</sup> (oku), 10<sup>12</sup> (chō), and 10<sup>16</sup> (kei).

In theory, these could be mapped to extended size units as follows:

```
man (万, 10⁴ bytes) → manbytes
oku (億, 10⁸ bytes) → okubytes
chō (兆, 10¹² bytes) → chobytes
kei (京, 10¹⁶ bytes) → keibytes
```

However, **miruzo currently uses manbytes only**.

For image loading and rendering strategies, a 10<sup>4</sup>-byte scale
provides sufficient granularity, and there is no practical need to introduce
higher-order units at this time.

These names are included solely to illustrate the consistency of the underlying
numeric system and to leave room for potential future exploration.

## 6. Non-goals

Manbytes is not intended for precise storage accounting or capacity
measurement. It is also not designed for strict bandwidth estimation, billing
calculations, or as an input representation for machine learning models.


## 7. This is cultural — but also practical

Manbytes is a unit inspired by the Japanese numeric system based on
10<sup>4</sup> groupings. However, its motivation is not cultural novelty.

In UI and UX implementations such as image rendering and loading strategies, it
provides a scale that is *well-matched to human perception*.

Manbytes is designed:

- not to represent precise storage size,
- but to classify asset weight in a way that supports decision-making.

In miruzo, this property is used to guide image display behavior and loading
heuristics.

It is a lightweight size scale, born at the intersection of culture and
engineering.


## 8. Summary

- Manbytes is a size unit based on 10<sup>4</sup> bytes
- Inspired by Japanese numeric grouping, but driven by practical needs
- Well-suited for *stable classification* of image sizes
- Used as a heuristic for UX and loading decisions
- miruzo standardizes on manbytes only
