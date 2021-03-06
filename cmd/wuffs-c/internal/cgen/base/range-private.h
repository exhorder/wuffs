// After editing this file, run "go generate" in the parent directory.

// Copyright 2017 The Wuffs Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ---------------- Ranges and Rects

static inline uint32_t  //
wuffs_base__range_ii_u32__get_min_incl(const wuffs_base__range_ii_u32* r) {
  return r->min_incl;
}

static inline uint32_t  //
wuffs_base__range_ii_u32__get_max_incl(const wuffs_base__range_ii_u32* r) {
  return r->max_incl;
}

static inline uint32_t  //
wuffs_base__range_ie_u32__get_min_incl(const wuffs_base__range_ie_u32* r) {
  return r->min_incl;
}

static inline uint32_t  //
wuffs_base__range_ie_u32__get_max_excl(const wuffs_base__range_ie_u32* r) {
  return r->max_excl;
}

static inline uint64_t  //
wuffs_base__range_ii_u64__get_min_incl(const wuffs_base__range_ii_u64* r) {
  return r->min_incl;
}

static inline uint64_t  //
wuffs_base__range_ii_u64__get_max_incl(const wuffs_base__range_ii_u64* r) {
  return r->max_incl;
}

static inline uint64_t  //
wuffs_base__range_ie_u64__get_min_incl(const wuffs_base__range_ie_u64* r) {
  return r->min_incl;
}

static inline uint64_t  //
wuffs_base__range_ie_u64__get_max_excl(const wuffs_base__range_ie_u64* r) {
  return r->max_excl;
}

// ---------------- Ranges and Rects (Utility)

static inline wuffs_base__range_ii_u32  //
wuffs_base__utility__make_range_ii_u32(const wuffs_base__utility* ignored,
                                       uint32_t min_incl,
                                       uint32_t max_incl) {
  return ((wuffs_base__range_ii_u32){
      .min_incl = min_incl,
      .max_incl = max_incl,
  });
}

static inline wuffs_base__range_ie_u32  //
wuffs_base__utility__make_range_ie_u32(const wuffs_base__utility* ignored,
                                       uint32_t min_incl,
                                       uint32_t max_excl) {
  return ((wuffs_base__range_ie_u32){
      .min_incl = min_incl,
      .max_excl = max_excl,
  });
}

static inline wuffs_base__range_ii_u64  //
wuffs_base__utility__make_range_ii_u64(const wuffs_base__utility* ignored,
                                       uint64_t min_incl,
                                       uint64_t max_incl) {
  return ((wuffs_base__range_ii_u64){
      .min_incl = min_incl,
      .max_incl = max_incl,
  });
}

static inline wuffs_base__range_ie_u64  //
wuffs_base__utility__make_range_ie_u64(const wuffs_base__utility* ignored,
                                       uint64_t min_incl,
                                       uint64_t max_excl) {
  return ((wuffs_base__range_ie_u64){
      .min_incl = min_incl,
      .max_excl = max_excl,
  });
}

static inline wuffs_base__rect_ii_u32  //
wuffs_base__utility__make_rect_ii_u32(const wuffs_base__utility* ignored,
                                      uint32_t min_incl_x,
                                      uint32_t min_incl_y,
                                      uint32_t max_incl_x,
                                      uint32_t max_incl_y) {
  return ((wuffs_base__rect_ii_u32){
      .min_incl_x = min_incl_x,
      .min_incl_y = min_incl_y,
      .max_incl_x = max_incl_x,
      .max_incl_y = max_incl_y,
  });
}

static inline wuffs_base__rect_ie_u32  //
wuffs_base__utility__make_rect_ie_u32(const wuffs_base__utility* ignored,
                                      uint32_t min_incl_x,
                                      uint32_t min_incl_y,
                                      uint32_t max_excl_x,
                                      uint32_t max_excl_y) {
  return ((wuffs_base__rect_ie_u32){
      .min_incl_x = min_incl_x,
      .min_incl_y = min_incl_y,
      .max_excl_x = max_excl_x,
      .max_excl_y = max_excl_y,
  });
}
