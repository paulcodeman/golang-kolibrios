#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

extern void* malloc(size_t size);

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

char* runtime_concatstrings(const char** strings, size_t count) {
    size_t total_length;
    size_t offset;
    size_t index;
    char* result;

    if (strings == NULL || count == 0) {
        return NULL;
    }

    total_length = 0;
    for (index = 0; index < count; index++) {
        if (strings[index] != NULL) {
            total_length += kos_strlen(strings[index]);
        }
    }

    result = (char*)malloc(total_length + 1);
    if (result == NULL) {
        return NULL;
    }

    offset = 0;
    for (index = 0; index < count; index++) {
        const char* current;
        size_t length;

        current = strings[index];
        if (current == NULL) {
            continue;
        }

        length = kos_strlen(current);
        kos_memcpy(result + offset, current, length);
        offset += length;
    }

    result[offset] = '\0';
    return result;
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

bool runtime_strequal(const char* left, const char* right) {
    if (left == NULL || right == NULL) {
        return false;
    }

    return kos_strcmp(left, right) == 0;
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
__asm__(".set runtime.writeBarrier, runtime_write_barrier");

__asm__(".global runtime.gcWriteBarrier");
__asm__(".set runtime.gcWriteBarrier, runtime_gc_write_barrier");

__asm__(".global runtime.strequal..f");
__asm__(".set runtime.strequal..f, runtime_strequal");
