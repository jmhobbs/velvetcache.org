#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if ! command -v pyftsubset &> /dev/null; then
  echo "pyftsubset not found, please install to continue"
  exit 1
fi

if ! [ -f "$__dir/SourceCodePro-Regular.ttf" ]; then
  echo "SourceCodePro-Regular.ttf not found, attempting to download..."
  curl \
    -L \
    -o "$__dir/SourceCodePro-Regular.ttf" \
    https://github.com/adobe-fonts/source-code-pro/raw/release/TTF/SourceCodePro-Regular.ttf
fi

echo "Removing old files..."

rm -f SourceCodePro-Regular-patched.ttf
rm -f SourceCodePro-Regular-subset*.woff2
rm -f source-code-pro-*.css

UNICODES=(
  "U+0020-007E" # ASCII Printables
  "U+00A1-00BF" # Latin-1 symbols
  "U+00D7"      # Multiplication sign
  "U+00F7"      # Division sign
  "U+2000-206F" # General Punctuation
  "U+2190-21FF" # Arrows
  "U+2200-22FF" # Mathematical Operators
  "U+2500-257F" # Box Drawing
)

unicodes_string="${UNICODES[*]}"

echo "Patching...."

python patch-font.py

echo "Subsetting..."

pyftsubset SourceCodePro-Regular-patched.ttf \
  --output-file="SourceCodePro-Regular-subset.woff2" \
  --flavor=woff2 \
  --no-hinting \
  --desubroutinize \
  --layout-features="ccmp" \
  --unicodes="${unicodes_string//${IFS:0:1}/,}"

SHA=$(sha256sum SourceCodePro-Regular-subset.woff2 | head -c 7)

mv SourceCodePro-Regular-subset.woff2 "SourceCodePro-Regular-subset-$SHA.woff2"

cat > "source-code-pro-$SHA.css" <<EOF
@font-face {
  font-family: 'Source Code Pro';
  font-style: normal;
  font-weight: 400;
  font-display: swap;
  src: url('SourceCodePro-Regular-subset-$SHA.woff2') format('woff2');
}
EOF

echo "-> SourceCodePro-Regular-subset-$SHA.woff2"
echo "-> source-code-pro-$SHA.css"
