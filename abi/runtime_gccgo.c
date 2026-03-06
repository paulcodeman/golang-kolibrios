#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

extern void* malloc(size_t size);
extern void* realloc(void* ptr, size_t size);
extern void free(void* ptr);

typedef struct {
    const char* str;
    intptr_t len;
} go_string;

static size_t kos_strlen(const char* str) {
    const char* cursor = str;
    while (*cursor != '\0') {
        cursor++;
    }
    return (size_t)(cursor - str);
}

static int kos_strcmp(const char* left, const char* right) {
    while (*left != '\0' && *left == *right) {
        left++;
        right++;
    }
    return (int)(*(const unsigned char*)left) - (int)(*(const unsigned char*)right);
}

static void* kos_memcpy(void* dest, const void* src, size_t size) {
    unsigned char* out = (unsigned char*)dest;
    const unsigned char* in = (const unsigned char*)src;

    while (size-- > 0) {
        *out++ = *in++;
    }

    return dest;
}

static int kos_memcmp(const void* left, const void* right, size_t size) {
    const unsigned char* left_bytes = (const unsigned char*)left;
    const unsigned char* right_bytes = (const unsigned char*)right;
    size_t index;

    for (index = 0; index < size; index++) {
        if (left_bytes[index] != right_bytes[index]) {
            return (int)left_bytes[index] - (int)right_bytes[index];
        }
    }

    return 0;
}

static int runtime_write_barrier_enabled = 0;
static char* runtime_window_title_buffer = NULL;
static size_t runtime_window_title_capacity = 0;

bool runtime_memequal32(const unsigned char* left, const unsigned char* right, size_t size) {
    size_t index;
    const uint32_t* left_words;
    const uint32_t* right_words;

    if (left == NULL || right == NULL || size % 4 != 0) {
        return false;
    }

    left_words = (const uint32_t*)left;
    right_words = (const uint32_t*)right;

    for (index = 0; index < size / 4; index++) {
        if (left_words[index] != right_words[index]) {
            return false;
        }
    }

    return true;
}

static const char* runtime_prepare_window_title_impl(uint32_t prefix, int use_prefix, const char* src, intptr_t len) {
    char* resized;
    size_t needed;
    size_t offset;

    if (src == NULL) {
        return NULL;
    }

    if (len < 0) {
        len = 0;
    }

    offset = use_prefix ? 1u : 0u;
    needed = offset + (size_t)len + 1;
    if (runtime_window_title_buffer == NULL || needed > runtime_window_title_capacity) {
        resized = (char*)realloc(runtime_window_title_buffer, needed);
        if (resized == NULL) {
            return runtime_window_title_buffer;
        }

        runtime_window_title_buffer = resized;
        runtime_window_title_capacity = needed;
    }

    if (use_prefix) {
        runtime_window_title_buffer[0] = (char)prefix;
    }

    if (len > 0) {
        kos_memcpy(runtime_window_title_buffer + offset, src, (size_t)len);
    }
    runtime_window_title_buffer[offset + (size_t)len] = '\0';

    return runtime_window_title_buffer;
}

const char* runtime_prepare_window_title(const char* src, intptr_t len) {
    return runtime_prepare_window_title_impl(0, 0, src, len);
}

const char* runtime_prepare_window_title_with_prefix(uint32_t prefix, const char* src, intptr_t len) {
    return runtime_prepare_window_title_impl(prefix, 1, src, len);
}

char* runtime_alloc_cstring(const char* src, intptr_t len) {
    char* out;

    if (src == NULL) {
        return NULL;
    }

    if (len < 0) {
        len = 0;
    }

    out = (char*)malloc((size_t)len + 1);
    if (out == NULL) {
        return NULL;
    }

    if (len > 0) {
        kos_memcpy(out, src, (size_t)len);
    }
    out[len] = '\0';

    return out;
}

void runtime_free_cstring(void* ptr) {
    if (ptr != NULL) {
        free(ptr);
    }
}

uint32_t runtime_pointer_value(void* ptr) {
    return (uint32_t)(uintptr_t)ptr;
}

bool runtime_memequal8(const unsigned char* left, const unsigned char* right, size_t size) {
    size_t index;

    if (left == NULL || right == NULL) {
        return false;
    }

    for (index = 0; index < size; index++) {
        if (left[index] != right[index]) {
            return false;
        }
    }

    return true;
}

go_string runtime_concatstrings(uintptr_t ignored, const go_string* strings, size_t count) {
    size_t total_length;
    size_t offset;
    size_t index;
    char* result;
    go_string out;

    (void)ignored;

    if (strings == NULL || count == 0) {
        out.str = NULL;
        out.len = 0;
        return out;
    }

    total_length = 0;
    for (index = 0; index < count; index++) {
        if (strings[index].str != NULL && strings[index].len > 0) {
            total_length += (size_t)strings[index].len;
        }
    }

    result = (char*)malloc(total_length + 1);
    if (result == NULL) {
        out.str = NULL;
        out.len = 0;
        return out;
    }

    offset = 0;
    for (index = 0; index < count; index++) {
        go_string current;
        size_t length;

        current = strings[index];
        if (current.str == NULL || current.len <= 0) {
            continue;
        }

        length = (size_t)current.len;
        kos_memcpy(result + offset, current.str, length);
        offset += length;
    }

    result[offset] = '\0';
    out.str = result;
    out.len = (intptr_t)offset;
    return out;
}

void runtime_set_byte_string(unsigned char* dest, const unsigned char* src, size_t size) {
    if (dest == NULL || src == NULL) {
        return;
    }

    kos_memcpy(dest, src, size);
}

void runtime_write_barrier(void** slot, void* ptr) {
    if (slot != NULL) {
        *slot = ptr;
    }
}

void runtime_gc_write_barrier(void** slot, void* ptr) {
    runtime_write_barrier(slot, ptr);
}

bool runtime_strequal(go_string left, go_string right) {
    size_t length;

    if (left.len != right.len) {
        return false;
    }

    if (left.str == right.str) {
        return true;
    }

    if (left.str == NULL || right.str == NULL) {
        return false;
    }

    length = (size_t)left.len;
    return kos_memcmp(left.str, right.str, length) == 0;
}

int memcmp(const void* left, const void* right, size_t size) {
    if (left == NULL || right == NULL) {
        return left == right ? 0 : (left == NULL ? -1 : 1);
    }

    return kos_memcmp(left, right, size);
}

void runtime_panicmem(void) {
    for (;;) {
    }
}

void* __unsafe_get_addr(void* base, size_t offset) {
    if (base == NULL) {
        return NULL;
    }

    return (void*)((unsigned char*)base + offset);
}

__asm__(".global runtime.memequal32..f");
__asm__(".set runtime.memequal32..f, runtime_memequal32");

__asm__(".global runtime.memequal8..f");
__asm__(".set runtime.memequal8..f, runtime_memequal8");

__asm__(".global runtime.memequal");
__asm__(".set runtime.memequal, runtime_memequal8");

__asm__(".global runtime.concatstrings");
__asm__(".set runtime.concatstrings, runtime_concatstrings");

__asm__(".global runtime.SetByteString");
__asm__(".set runtime.SetByteString, runtime_set_byte_string");

__asm__(".global runtime.writeBarrier");
__asm__(".set runtime.writeBarrier, runtime_write_barrier_enabled");

__asm__(".global runtime.gcWriteBarrier");
__asm__(".set runtime.gcWriteBarrier, runtime_gc_write_barrier");

__asm__(".global runtime.strequal..f");
__asm__(".set runtime.strequal..f, runtime_strequal");

__asm__(".global runtime.panicmem");
__asm__(".set runtime.panicmem, runtime_panicmem");
