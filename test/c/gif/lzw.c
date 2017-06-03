// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

/*
To manually run this test:

for cc in clang gcc; do
  $cc -std=c99 -Wall -Werror lzw.c && ./a.out
  rm -f a.out
done

Each edition should print "PASS", amongst other information, and exit(0).
*/

#include "../../../gen/c/gif.c"
#include "../testlib/testlib.c"

const char* test_filename = "gif/lzw.c";

#define BUFFER_SIZE (1024 * 1024)

uint8_t global_dst_buffer[BUFFER_SIZE];
uint8_t global_src_buffer[BUFFER_SIZE];

void test_lzw_decode() {
  puffs_gif_lzw_decoder dec;
  puffs_gif_lzw_decoder_constructor(&dec, PUFFS_VERSION, 0);

  puffs_base_buf1 dst = {.ptr = global_dst_buffer, .cap = BUFFER_SIZE};
  puffs_base_buf1 src = {.ptr = global_src_buffer, .cap = BUFFER_SIZE};
  global_src_buffer[src.wi++] = 0x80;
  global_src_buffer[src.wi++] = 0x81;

  puffs_gif_status status =
      puffs_gif_lzw_decoder_decode(&dec, &dst, &src, false);
  if (status != puffs_gif_status_ok) {
    FAIL("test_lzw_decode: status: got %d, want %d", status,
         puffs_gif_status_ok);
    goto cleanup0;
  }
cleanup0:
  puffs_gif_lzw_decoder_destructor(&dec);
}

// The empty comments forces clang-format to place one element per line.
test tests[] = {
    test_lzw_decode,  //
    NULL,             //
};