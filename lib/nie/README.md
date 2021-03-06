# Naïve Image Formats: NIE, NII, NIA

Status: Draft (as of December 2018). There is no compatibility guarantee yet.


## NIE: Still Images

NIE is an easily parsed, uncompressed, lossless format for still (single frame)
images. The 16 byte header:

- 4 bytes of 'magic': \[0x6E, 0xC3, 0xAF, 0x45\], the UTF-8 encoding of "nïE".
- 4 bytes of version-and-configuration: \[0xFF, 0x62 or 0x72, 0x6E or 0x70,
  0x34 or 0x38\].
  - The first byte denotes the overall NIE/NII/NIA format version. 0xFF (which
	is not valid UTF-8) denotes version 1. There are no other valid versions at
	this time.
  - The second byte, either an ASCII 'b' or an ASCII 'r', denotes whether the
	pixel's data is in BGRA or RGBA order (in memory order, independent of CPU
	endianness).
  - The third byte, either an ASCII 'n' or an ASCII 'p', denotes whether the
	payload contains non-premultiplied or premultiplied alpha.
  - The fourth byte, either an ASCII '4' or an ASCII '8', denotes whether there
	are 4 or 8 bytes per pixel.
  - Future format versions may allow other byte values, but in version 1, it
	must be '\xFF', then 'b' or 'r', then 'n' or 'p', then '4' or '8'.
- 4 bytes little-endian `uint32` width.
- 4 bytes little-endian `uint32` height.

The payload:

- 4 or 8 bytes per pixel. W×H pixels in row-major order. Row-major means that
  horizontally adjacent pixels are adjacent in memory. Values are
  little-endian. For example, with BGRA order and 8 bytes per pixel, those
  bytes are \[B₀, B₁, G₀, G₁, R₀, R₁, A₀, A₁\]. The ₀ and ₁ subscripts denote
  the low and high bytes of the `uint16`.


## NII: Animated Images, Timing Index Only, Out-of-Band Frames

NII is an index for animated (multiple frame) images. In video compression
terminology, every NII frame is an I-frame, also known as a keyframe. A NII
file doesn't contain the images per se, only the duration that each frame
should be shown.

The per-frame images, not part of a NII file, may be NIE images, but they may
also be in other formats, such as PNG or WebP or a heterogenous mixture. They
may be static files (possibly with systematic filenames such as
`frame000000.png`, `frame000001.png`, etc.) or dynamically generated. That is
for each application to decide, and out of scope of this specification.


### NII Header

The 32 byte NII header:

- 4 bytes of 'magic': \[0x6E, 0xC3, 0xAF, 0x49\]. The final byte differs from
  NIE: an ASCII 'I' instead of an ASCII 'E'.
- 4 bytes of version-and-padding, all 0xFF.
- 4 bytes little-endian `uint32` width.
- 4 bytes little-endian `uint32` height.
- 4 bytes little-endian `uint32` LCValue. See below.
- 4 bytes little-endian `uint32` NFValue. See below.
- 8 bytes little-endian `int64` Overall-TValue. See below.

A zero LCValue means that the animation loops forever. Non-zero means a loop
count: the animation is played LCValue times. This is the
[APNG](https://wiki.mozilla.org/APNG_Specification) meaning, not the
[GIF](https://www.w3.org/Graphics/GIF/spec-gif89a.txt) meaning (the number of
times to repeat the loop _after_ the first play). The two meanings differ by 1.

The Overall-TValue must not be positive. If it is zero, the NFValue must also
be zero. If it is negative, it must equal the final frame's TValue (the final 8
bytes of the NII file, see the NII Payload section below), or if there are no
frames, it must equal -1.

A negative Overall-TValue means the overall animation length, both in terms of
time and in terms of number of frames, is known at the time the 32 byte NII
header was written. The overall time is the OCDD (see the NII Payload section
below). The overall number of frames is the NFValue. Animations lasting longer
than `((1<<32)-1)` frames, more than 4 billion frames, are not representable in
the NII format.

A zero Overall-TValue means that the animation has at least one frame, but the
overall animation length is not known until the end of the NII file is reached.
For example, a one-pass GIF to NII converter might not know the total duration
or the total number of frames until the entire GIF file is processed, at which
point the NII header may already have been written.

To repeat, if an animation has no frames, its NFValue and Overall-TValue must
be 0 and -1 respectively. If an animation has at least one frame, its NFValue
and Overall-TValue pair must be either (zero, zero) or (positive, negative).


### NII Payload

The payload is a sequence of 0 or more frames, exactly 8 bytes per frame:

- 8 bytes little-endian `int64` timing value (TValue).

Every TValue must be non-negative, except for the final TValue, which must be
negative. A frame's CDD (cumulative display duration) equals its TValue for all
but the final frame, and equals the bitwise-complement (or equivalently, "xor
with -1") of the TValue for that final frame. CDD values must therefore be
non-negative int64 values.

The OCDD (overall cumulative display duration) is the final frame's CDD. If the
animation contains no frames, then the OCDD is 0 and the Overall-TValue is -1.

For example, if an animation has four frames, to be displayed for 1 second, 2
seconds, 0 seconds and finally 4.5 seconds, then the CDD's are 1s, 3s, 3s and
7.5s. NII's unit of time is [flicks](https://github.com/OculusVR/Flicks): one
flick (frame-tick) is 1 / 705\_600\_000 of a second. Continuing our example,
the CDDs (in decimal and then hexadecimal) are:

- 705\_600\_000 × 1.0 = 0\_705\_600\_000 = 0x0000\_0000\_2A0E\_9A00.
- 705\_600\_000 × 3.0 = 2\_116\_800\_000 = 0x0000\_0000\_7E2B\_CE00.
- 705\_600\_000 × 3.0 = 2\_116\_800\_000 = 0x0000\_0000\_7E2B\_CE00.
- 705\_600\_000 × 7.5 = 5\_292\_000\_000 = 0x0000\_0001\_3B6D\_8300.

The TValue bytes (little-endian `int64` encoded) are therefore:

- 0x00 0x9A 0x0E 0x2A 0x00 0x00 0x00 0x00
- 0x00 0xCE 0x2B 0x7E 0x00 0x00 0x00 0x00
- 0x00 0xCE 0x2B 0x7E 0x00 0x00 0x00 0x00
- 0xFF 0x7C 0x92 0xC4 0xFE 0xFF 0xFF 0xFF

It is invalid for one frame's CDD value to be less than the previous frame's
CDD value. Animations lasting longer than `((1<<63)-1)` flicks, more than 400
years, are not representable in the NII format.


## NIA: Animated Images, In-band Frames

NIA is like an NII file where the per-frame still images are NIE files
interleaved between the NII TValues.

The NIA header is the same as the 32 byte NII header, except that the 4 byte
'magic' ends in an ASCII 'A' instead of an ASCII 'I', and the 5th to 8th bytes
are version-and-configuration (the same as for NIE), instead of NII's
version-and-padding. The range of valid version-and-configuration bytes is the
same for NIA as it is for NIE.

The payload is a sequence of 0 or more frames. Each frame is:

- A complete NIE image: header and payload. The outer NIA and inner NIE
  must have the same 12 bytes of version-and-configuration, width and height.
- Either 0 or 4 bytes of padding. If present, it must be all zeroes. The
  padding ensures that the size of the padded NIE image a multiple of 8 bytes,
  so that every TValue field is 8 byte aligned. The padding size is 4 if and
  only if there are 4 (not 8) bytes per pixel and both the width and height are
  odd. A C programming language expression for its presence is
  `((bytes_per_pixel == 4) && (width & height & 1))`.
- 8 bytes little-endian int64 TValue, the same meaning and constraints as NII.


# Commentary


## Motivation

One motivating example is securely decoding untrusted images, perhaps uploaded
from potentially malicious actors. Codec libraries have been a rich source of
software security vulnerabilities in the past. One response is to split off
such code into a separate, sandboxed process that reads the compressed image
and writes the equivalent NIE/NIA image, perhaps through pipes or shared
memory. The untrusted codec library processes the untrusted data within the
sandboxed worker process. The unsandboxed manager process only needs to handle
the much simpler NIE/NIA format. That format can be further simplified by the
manager mandating a fixed version-and-configuration, such as v1-"bn4".

Another example is connecting a series of independent image manipulation
programs, each component reading, transforming and then writing a NIE/NIA
image. Such filters can be written in simple programming languages and
connected with Unix-style pipes.

Another example is storing 'golden images' for codec development. Given a
corpus of test images in a compressed format (e.g. a corpus of PNG files), it
is useful to store their expected decodings for comparison, but those golden
test files should be encoded in an alternative format, such as NIE/NIA.


## Magic

The 4 bytes of 'magic' are the UTF-8 encoding of the non-ASCII strings "nïE",
"nïI" and "nïA". The unusual capitalization lessens the chance of plain text
data accidentally matching these magic bytes.


## Alpha Premultiplication

For premultiplied alpha, it is valid for a pixel's blue, green or red values to
be greater than its alpha value. Interpretation of such super-saturated colors
is out of scope of this specification.

A program that simply extracts a subset of a NIA's frames as a new NIA
animation is not required to examine or re-encode every payload byte in order
to always output valid NIA data.


## Random Access

Given a NIA animation's bytes per pixel, width and height _B_, _W_ and _H_, the
offset and length of the _i_'th frame's NIE data within that NIA is a simple
computation (but remember to check for overflow):

- length = roundup8(_B_ × _W_ × _H_) + 16
- offset = ((length + 8) × _i_) + 32

The roundup8 function rounds its argument up to the nearest multiple of 8.

This is random access by frame index (the "_i_" in "the _i_'th frame"), not by
time, as different frames can have different display durations.


## Cumulative Display Duration

Other animation formats, like
[APNG](https://wiki.mozilla.org/APNG_Specification) and
[GIF](https://www.w3.org/Graphics/GIF/spec-gif89a.txt), provide display
durations relative to the previous frame, not relative to the initial frame.
The two schemes are equivalent, in that from a complete stream, either one can
be derived from the other. NII / NIA frames report the cumulative number so
that random access by time can be implemented as a binary search, given random
access by frame index.

Suppose we are given a time _t_ ≥ 0 and want to find the frame to show at that
time. First, there may be no such frame, if the animation contains no frames.

Otherwise, let _o_ be the OCDD, so that _o_ ≥ 0. If _o_ is zero, the frame to
show is the final frame of the animation, and no further computation is
necessary.

Otherwise, calculate the number of loops that would complete by time t: _n_ =
_t_ / _o_, rounding down to the nearest integer. If the LCValue is non-zero (as
zero means loop forever) and _n_ ≥ LCValue then the frame to show is the final
frame.

Otherwise, calculate _t′_ = _t_ - (_n_ × _o_), the time 'modulo' _o_. Binary
search to find the smallest _i_ ≥ 0 such that both CDD(_i_) > _t′_ and the
_i_th frame is non-instantaneous. CDD(_i_) is the cumulative display duration
for frame _i_. The first frame is instantaneous if its CDD is zero. Any other
frame is instantaneous if its CDD equals its previous frame's CDD.


## Integer Arithmetic Overflow

Parsing NIE data is almost trivial, but care should be taken to avoid integer
arithmetic overflow when calculating the pixel buffer size from fields in the
NIE header. For example, a C programming language statement like `size_t
row_size = bytes_per_pixel * width;` is incorrect without additional prior
checks. A careful C implementation is:

```c
#include <stdbool.h>
#include <stdint.h>

// nie_payload_size calculates the size in bytes of a NIE payload, given the
// metadata from the NIE header: bytes_per_pixel, width and height. The max
// argument, not defined in the metadata, is the caller's maximum acceptable
// payload size. For example, pass SIZE_MAX for max, or pass a smaller value if
// you wish to limit the memory required to decode an arbitrary NIE file (and
// reject otherwise valid NIE files that would require more memory).
//
// That size is essentially (bytes_per_pixel * width * height), but this
// function checks for integer arithmetic overflow. It also checks that the
// result pointer is non-NULL, the result (the calculated payload size) is less
// than or equal to max, and that the bytes_per_pixel is either 4 or 8.
//
// The bool return value is whether all checks pass. On success, it sets
// *result to the payload size.
bool nie_payload_size(size_t* result,
                      size_t max,
                      uint32_t bytes_per_pixel,
                      uint32_t width,
                      uint32_t height) {
  if ((result == NULL) || ((bytes_per_pixel != 4) && (bytes_per_pixel != 8))) {
    return false;
  }
  uint64_t n = ((uint64_t)width) * ((uint64_t)height);

  // bpp_shift is 2 or 3, depending on bytes_per_pixel being 4 or 8.
  uint32_t bpp_shift = 2 + (bytes_per_pixel >> 3);
  if (n > (max >> bpp_shift)) {
    return false;
  }
  n <<= bpp_shift;

  *result = (size_t)n;
  return true;
}
```


## Color Management and Other Metadata

There is no facility for describing color spaces, gamma, palettes or other
metadata such as EXIF information. For example, when using a sandboxed worker
process to convert from a PNG image (with an embedded color profile) to a NIE
image, the target color space should be provided to the worker out-of-band.


## YUV Color

There is no facility for explicitly describing YUV or Y'CbCr color. Converting
between NIE/NIA and formats such as
[JPEG](http://www.w3.org/Graphics/JPEG/itu-t81.pdf) or [WebP
Lossy](https://developers.google.com/speed/webp/docs/riff_container) is a lossy
process, although JPEG and WebP Lossy are lossy formats to begin with.


## Bytes versus Octets

It was not always the case, historically, but in this specification, `byte` is
synonymous with `octet` and `uint8`. Similarly, signed integers are assumed to
use two's complement representation.


# Filename Extensions and MIME Types

The recommended filename extensions are `.nie`, `.nii` and `.nia`.

The recommended MIME types are `image/nie`, `image/nii` and `image/nia`.


## Why NIE and not NIF (for Naïve Image Format)?

The `.nif` filename extension is already used by the NetImmerse / Gamebryo game
engine. Instead, you can think of `.nie` as derived from the word "naïve".


## Pronunciation

I pronounce "NIE", "NII" and "NIA" as "naɪ'i", "naɪ'aɪ"  and "naɪ'eɪ", ending
in a long "E", "I" or "A" sound. It's definitely a hard "N", not a soft one.


# Related Work

In general, being uncompressed, NIE/NIA is not intended for long term storage
or for delivery over a network. Alternative lossless image formats like
[PNG](https://www.w3.org/TR/2003/REC-PNG-20031110/) or [WebP
Lossless](https://developers.google.com/speed/webp/docs/riff_container) are far
more suitable.

Amongst uncompressed formats, [Farbfeld](https://tools.suckless.org/farbfeld/)
is a very similar design. The main difference is that NIE is little-endian (the
same as the vast majority of CPUs) where Farbfeld is big-endian. NIE/NII/NIA
also allows animated images, not just still images, and is also configurable
with respect to BGRA versus RGBA order, non-premultiplied versus premultiplied
alpha, and 8 versus 16 bit channels. Configurability is admittedly a trade-off:
one person's configuration parameter is another person's unnecessary bloat.
We're not saying that Farbfeld is a bad design, just a different design that
has chosen different trade-offs.

[PAM](http://netpbm.sourceforge.net/doc/pam.html) (Portable Arbitrary Map) and
its related formats (PBM, PGM, PPM) are also very similar. The main difference
is that their headers are ASCII text and of variable length. Such headers are
easier for humans to read but harder for computers to parse. Again, NIE/NII/NIA
also allows animated images. Some programs process concatenated streams of PAM
images as an animation, but that doesn't contain explicit timing information.

[VIPS](https://jcupitt.github.io/libvips/API/current/file-format.html) is the
intermediate format for the [libvips](https://libvips.github.io/libvips/) image
processing library. Pixel data is stored in a straightforward fashion, but
metadata is stored as XML, which is not simple, and is a potential security
concern.

Amongst flexible formats,
[BMP](https://www.fileformat.info/format/bmp/egff.htm),
[PNG](https://www.w3.org/TR/2003/REC-PNG-20031110/) and its animated variant
[APNG](https://wiki.mozilla.org/APNG_Specification),
[Targa](http://tfc.duke.free.fr/coding/tga_specs.pdf) and
[TIFF](http://partners.adobe.com/public/developer/en/tiff/TIFF6.pdf) all have
their built-in compression being optional. Some of those uncompressed forms
have been used as simple interchange formats. Compared to those, NIE has fewer
features, which is obviously unfavorable if you want such features, but can be
favorable if you prefer simplicity and code auditability. A NIE file is still
simpler than an uncompressed BMP file and a NIA file is still simpler than an
uncompressed APNG file.

For example, when using a sandboxed worker process to convert to an
intermediate format, we don't want to use "all possible valid TIFF images" as
the intermediate format, because a compromised worker process could generate a
malicious intermediate image that provoked a security bug in an unsandboxed
TIFF decoder. Instead, we would want only a very simple subset of TIFF, but it
can be easier (in terms of implementation and in terms of a code audit) to
write an (almost trivial) NIE parser from scratch instead of whittling down a
full-featured TIFF library and hoping that you've eliminated all the potential
security vulnerabilities.


---

Updated on December 2018.
