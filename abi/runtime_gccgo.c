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

typedef struct {
    unsigned char* values;
    intptr_t len;
    intptr_t cap;
} go_slice;

typedef struct {
    uintptr_t size;
    uintptr_t ptrdata;
    uint32_t hash;
    uint8_t tflag;
    uint8_t align;
    uint8_t field_align;
    uint8_t kind;
    bool (**equal)(const void* left, const void* right);
    const void* gcdata;
    const go_string* name;
    const void* uncommon;
    const void* ptr_to_this;
} go_type_descriptor;

typedef struct {
    const go_type_descriptor* type;
} go_interface_method_table;

typedef struct {
    const go_interface_method_table* methods;
    const void* data;
} go_interface;

typedef struct {
    const go_type_descriptor* type;
    const void* data;
} go_empty_interface;

typedef struct {
    const go_string* name;
    const go_string* package_path;
    const void* methods;
    uint32_t method_count;
    uint32_t exported_method_count;
} go_uncommon_type;

typedef struct {
    const go_type_descriptor common;
    const void* methods;
    uint32_t method_count;
    uint32_t exported_method_count;
} go_interface_type_descriptor;

typedef struct {
    const go_string* name;
    const go_string* package_path;
    const go_type_descriptor* type;
} go_interface_method_descriptor;

typedef struct {
    const go_string* name;
    const go_string* package_path;
    const go_type_descriptor* interface_type;
    const go_type_descriptor* concrete_type;
    void* function;
} go_named_type_method_descriptor;

typedef struct {
    go_interface value;
    bool ok;
} go_interface_assert_result;

typedef bool (*go_equal_function)(const void* left, const void* right);

#define GO_TYPE_KIND_DIRECT_IFACE 0x20u
#define GO_TYPE_KIND_MASK 0x1Fu
#define GO_TYPE_KIND_INTERFACE 0x14u

typedef struct {
    uintptr_t size;
} go_type_size_only_descriptor;

void runtime_panicmem(void);
void runtime_typedmemmove(const go_type_descriptor* descriptor, void* dest, const void* src);

static const char runtime_hex_digits[] = "0123456789ABCDEF";

static int kos_memcmp(const void* left, const void* right, size_t size);

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

static bool runtime_string_equals(const go_string* left, const go_string* right) {
    size_t size;

    if (left == right) {
        return true;
    }
    if (left == NULL || right == NULL) {
        return false;
    }
    if (left->len != right->len) {
        return false;
    }
    if (left->len == 0) {
        return true;
    }
    if (left->str == NULL || right->str == NULL) {
        return false;
    }

    size = (size_t)left->len;
    return kos_memcmp(left->str, right->str, size) == 0;
}

static void* kos_memcpy(void* dest, const void* src, size_t size) {
    unsigned char* out = (unsigned char*)dest;
    const unsigned char* in = (const unsigned char*)src;

    while (size-- > 0) {
        *out++ = *in++;
    }

    return dest;
}

static void* kos_memmove(void* dest, const void* src, size_t size) {
    unsigned char* out;
    const unsigned char* in;
    size_t index;

    if (dest == src || size == 0) {
        return dest;
    }

    out = (unsigned char*)dest;
    in = (const unsigned char*)src;
    if (out < in || out >= in + size) {
        return kos_memcpy(dest, src, size);
    }

    for (index = size; index > 0; index--) {
        out[index - 1] = in[index - 1];
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

static void* kos_memset(void* dest, int value, size_t size) {
    unsigned char* out = (unsigned char*)dest;

    while (size-- > 0) {
        *out++ = (unsigned char)value;
    }

    return dest;
}

static void runtime_debug_write_byte(unsigned char value) {
    uint32_t eax;
    uint32_t ebx;
    uint32_t ecx;

    eax = 63;
    ebx = 1;
    ecx = (uint32_t)value;
    __asm__ volatile("int $0x40"
                     : "+a"(eax), "+b"(ebx), "+c"(ecx)
                     :
                     : "edx", "esi", "edi", "memory", "cc");
}

static void runtime_debug_write_bytes(const char* value, size_t size) {
    size_t index;

    if (value == NULL) {
        return;
    }

    for (index = 0; index < size; index++) {
        runtime_debug_write_byte((unsigned char)value[index]);
    }
}

static void runtime_debug_write_cstring(const char* value) {
    if (value == NULL) {
        return;
    }

    runtime_debug_write_bytes(value, kos_strlen(value));
}

static void runtime_debug_write_hex32(uint32_t value) {
    int shift;

    runtime_debug_write_cstring("0x");
    for (shift = 28; shift >= 0; shift -= 4) {
        runtime_debug_write_byte((unsigned char)runtime_hex_digits[(value >> shift) & 0x0F]);
    }
}

static void runtime_debug_write_newline(void) {
    runtime_debug_write_byte('\r');
    runtime_debug_write_byte('\n');
}

__attribute__((noreturn)) static void runtime_exit_process(void) {
    int32_t eax;

    eax = -1;
    __asm__ volatile("int $0x40"
                     : "+a"(eax)
                     :
                     : "ebx", "ecx", "edx", "esi", "edi", "memory", "cc");
    for (;;) {
    }
}

__attribute__((noreturn)) static void runtime_fail_simple(const char* reason) {
    runtime_debug_write_cstring("runtime panic: ");
    runtime_debug_write_cstring(reason);
    runtime_debug_write_newline();
    runtime_exit_process();
}

__attribute__((noreturn)) static void runtime_fail_pair(const char* reason, const char* first_name, uint32_t first_value, const char* second_name, uint32_t second_value) {
    runtime_debug_write_cstring("runtime panic: ");
    runtime_debug_write_cstring(reason);
    runtime_debug_write_cstring(" (");
    runtime_debug_write_cstring(first_name);
    runtime_debug_write_cstring("=");
    runtime_debug_write_hex32(first_value);
    runtime_debug_write_cstring(", ");
    runtime_debug_write_cstring(second_name);
    runtime_debug_write_cstring("=");
    runtime_debug_write_hex32(second_value);
    runtime_debug_write_cstring(")");
    runtime_debug_write_newline();
    runtime_exit_process();
}

static size_t kos_slice_allocation_size(const go_type_descriptor* descriptor, intptr_t count) {
    size_t element_size;

    if (count < 0) {
        runtime_panicmem();
    }

    if (count == 0) {
        return 0;
    }

    element_size = 0;
    if (descriptor != NULL) {
        element_size = (size_t)descriptor->size;
    }

    if (element_size == 0) {
        return 1;
    }

    if ((size_t)count > ((size_t)-1) / element_size) {
        runtime_panicmem();
    }

    return (size_t)count * element_size;
}

static int runtime_write_barrier_enabled = 0;
static char* runtime_window_title_buffer = NULL;
static size_t runtime_window_title_capacity = 0;

static bool runtime_memequal8_impl(const void* left, const void* right) {
    const unsigned char* left_bytes;
    const unsigned char* right_bytes;

    if (left == NULL || right == NULL) {
        return false;
    }

    left_bytes = (const unsigned char*)left;
    right_bytes = (const unsigned char*)right;
    return left_bytes[0] == right_bytes[0];
}

static bool runtime_memequal16_impl(const void* left, const void* right) {
    const uint16_t* left_words;
    const uint16_t* right_words;

    if (left == NULL || right == NULL) {
        return false;
    }

    left_words = (const uint16_t*)left;
    right_words = (const uint16_t*)right;
    return left_words[0] == right_words[0];
}

static bool runtime_memequal32_impl(const void* left, const void* right) {
    const uint32_t* left_words;
    const uint32_t* right_words;

    if (left == NULL || right == NULL) {
        return false;
    }

    left_words = (const uint32_t*)left;
    right_words = (const uint32_t*)right;
    return left_words[0] == right_words[0];
}

static bool runtime_memequal64_impl(const void* left, const void* right) {
    const uint32_t* left_words;
    const uint32_t* right_words;

    if (left == NULL || right == NULL) {
        return false;
    }

    left_words = (const uint32_t*)left;
    right_words = (const uint32_t*)right;
    return left_words[0] == right_words[0] &&
           left_words[1] == right_words[1];
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

static bool runtime_memequal_impl(const void* left, const void* right, size_t size) {
    if (left == NULL || right == NULL) {
        return false;
    }

    return kos_memcmp(left, right, size) == 0;
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

void* runtime_makeslice(const go_type_descriptor* descriptor, intptr_t len, intptr_t cap) {
    size_t total_size;
    void* memory;

    if (len < 0 || cap < len) {
        runtime_panicmem();
    }

    if (cap == 0) {
        return NULL;
    }

    total_size = kos_slice_allocation_size(descriptor, cap);

    memory = malloc(total_size);
    if (memory == NULL) {
        return NULL;
    }

    kos_memset(memory, 0, total_size);
    return memory;
}

go_slice runtime_growslice(const go_type_descriptor* descriptor, void* old_values, intptr_t old_len, intptr_t old_cap, intptr_t new_len) {
    go_slice result;
    size_t old_size;
    size_t new_size;
    intptr_t new_cap;
    void* memory;

    result.values = NULL;
    result.len = 0;
    result.cap = 0;

    if (old_len < 0 || old_cap < old_len || new_len < old_len) {
        runtime_panicmem();
    }

    new_cap = old_cap;
    if (new_cap < 1) {
        new_cap = 1;
    }

    while (new_cap < new_len) {
        if (new_cap > INTPTR_MAX / 2) {
            new_cap = new_len;
            break;
        }
        new_cap *= 2;
    }
    if (new_cap < new_len) {
        new_cap = new_len;
    }

    new_size = kos_slice_allocation_size(descriptor, new_cap);
    memory = malloc(new_size);
    if (memory == NULL) {
        return result;
    }

    kos_memset(memory, 0, new_size);
    old_size = kos_slice_allocation_size(descriptor, old_len);
    if (old_values != NULL && old_size > 0) {
        kos_memmove(memory, old_values, old_size);
    }

    result.values = (unsigned char*)memory;
    result.len = new_len;
    result.cap = new_cap;
    return result;
}

void runtime_typedmemmove(const go_type_descriptor* descriptor, void* dest, const void* src) {
    size_t size;

    if (dest == NULL || src == NULL || dest == src) {
        return;
    }

    size = 0;
    if (descriptor != NULL) {
        size = (size_t)descriptor->size;
    }

    if (size == 0) {
        return;
    }

    kos_memmove(dest, src, size);
}

go_string runtime_slicebytetostring(void* ignored, const unsigned char* src, intptr_t len) {
    char* out;
    go_string result;

    (void)ignored;

    if (src == NULL || len <= 0) {
        result.str = NULL;
        result.len = 0;
        return result;
    }

    out = (char*)malloc((size_t)len + 1);
    if (out == NULL) {
        result.str = NULL;
        result.len = 0;
        return result;
    }

    kos_memcpy(out, src, (size_t)len);
    out[len] = '\0';

    result.str = out;
    result.len = len;
    return result;
}

go_slice runtime_stringtoslicebyte(void* ignored, const char* src, intptr_t len) {
    go_slice result;

    (void)ignored;

    result.values = NULL;
    result.len = 0;
    result.cap = 0;
    if (src == NULL || len <= 0) {
        return result;
    }

    result.values = (unsigned char*)malloc((size_t)len);
    if (result.values == NULL) {
        return result;
    }

    kos_memcpy(result.values, src, (size_t)len);
    result.len = len;
    result.cap = len;
    return result;
}

void runtime_write_barrier(void** slot, void* ptr) {
    if (slot != NULL) {
        *slot = ptr;
    }
}

void runtime_gc_write_barrier(void** slot, void* ptr) {
    runtime_write_barrier(slot, ptr);
}

static bool runtime_strequal_impl(const void* left_value, const void* right_value) {
    const go_string* left;
    const go_string* right;
    size_t length;

    if (left_value == NULL || right_value == NULL) {
        return false;
    }

    left = (const go_string*)left_value;
    right = (const go_string*)right_value;

    if (left->len != right->len) {
        return false;
    }

    if (left->str == right->str) {
        return true;
    }

    if (left->str == NULL || right->str == NULL) {
        return false;
    }

    length = (size_t)left->len;
    return kos_memcmp(left->str, right->str, length) == 0;
}

static go_equal_function runtime_resolve_equal_function(const go_type_descriptor* descriptor) {
    if (descriptor == NULL || descriptor->equal == NULL) {
        return NULL;
    }

    return *descriptor->equal;
}

static bool runtime_type_descriptor_matches(const go_type_descriptor* left, const go_type_descriptor* right) {
    if (left == right) {
        return true;
    }
    if (left == NULL || right == NULL) {
        return false;
    }

    return left->size == right->size &&
        left->ptrdata == right->ptrdata &&
        left->hash == right->hash &&
        left->align == right->align &&
        left->field_align == right->field_align &&
        left->kind == right->kind &&
        runtime_string_equals(left->name, right->name);
}

static const go_named_type_method_descriptor* runtime_find_named_method(const go_uncommon_type* uncommon, const go_interface_method_descriptor* target_method) {
    const go_named_type_method_descriptor* methods;
    const go_named_type_method_descriptor* current;
    uint32_t index;

    if (uncommon == NULL || target_method == NULL || uncommon->methods == NULL || uncommon->method_count == 0) {
        return NULL;
    }

    methods = (const go_named_type_method_descriptor*)uncommon->methods;
    for (index = 0; index < uncommon->method_count; index++) {
        current = methods + index;
        if (!runtime_string_equals(current->name, target_method->name)) {
            continue;
        }
        if (!runtime_string_equals(current->package_path, target_method->package_path)) {
            continue;
        }
        if (!runtime_type_descriptor_matches(current->interface_type, target_method->type)) {
            continue;
        }

        return current;
    }

    return NULL;
}

static go_interface_method_table* runtime_build_interface_method_table(const go_interface_type_descriptor* target_interface, const go_type_descriptor* source_type) {
    const go_interface_method_descriptor* target_methods;
    const go_named_type_method_descriptor* source_method;
    const go_uncommon_type* uncommon;
    uintptr_t size;
    uintptr_t index;
    void** table_entries;

    if (target_interface == NULL || source_type == NULL) {
        return NULL;
    }
    if ((target_interface->common.kind & GO_TYPE_KIND_MASK) != GO_TYPE_KIND_INTERFACE) {
        return NULL;
    }

    size = sizeof(void*) + (uintptr_t)target_interface->method_count * sizeof(void*);
    table_entries = (void**)malloc((size_t)size);
    if (table_entries == NULL) {
        return NULL;
    }

    table_entries[0] = (void*)source_type;
    if (target_interface->method_count == 0 || target_interface->methods == NULL) {
        return (go_interface_method_table*)table_entries;
    }

    uncommon = (const go_uncommon_type*)source_type->uncommon;
    target_methods = (const go_interface_method_descriptor*)target_interface->methods;
    for (index = 0; index < (uintptr_t)target_interface->method_count; index++) {
        source_method = runtime_find_named_method(uncommon, target_methods + index);
        if (source_method == NULL || source_method->function == NULL) {
            free(table_entries);
            return NULL;
        }

        table_entries[index + 1] = source_method->function;
    }

    return (go_interface_method_table*)table_entries;
}

static void runtime_zero_typed_value(const go_type_descriptor* descriptor, void* dest) {
    size_t size;

    if (descriptor == NULL || dest == NULL) {
        return;
    }

    size = (size_t)descriptor->size;
    if (size == 0) {
        return;
    }

    kos_memset(dest, 0, size);
}

static void runtime_copy_typed_value(const go_type_descriptor* descriptor, void* dest, const void* src) {
    uintptr_t direct_value;
    size_t size;

    if (descriptor == NULL || dest == NULL) {
        return;
    }

    size = (size_t)descriptor->size;
    if (size == 0) {
        return;
    }

    if ((descriptor->kind & GO_TYPE_KIND_DIRECT_IFACE) != 0) {
        direct_value = (uintptr_t)src;
        kos_memcpy(dest, &direct_value, size);
        return;
    }

    if (src == NULL) {
        runtime_zero_typed_value(descriptor, dest);
        return;
    }

    runtime_typedmemmove(descriptor, dest, src);
}

static bool runtime_value_equal(const go_type_descriptor* descriptor, const void* left_data, const void* right_data) {
    go_equal_function equal;

    if (descriptor == NULL) {
        return true;
    }

    if ((descriptor->kind & GO_TYPE_KIND_DIRECT_IFACE) != 0) {
        return left_data == right_data;
    }

    equal = runtime_resolve_equal_function(descriptor);
    if (equal == NULL) {
        runtime_fail_simple("equality on non-comparable type");
    }

    return equal(left_data, right_data);
}

bool runtime_efaceeq(const go_type_descriptor* left_type, const void* left_data, const go_type_descriptor* right_type, const void* right_data) {
    if (left_type != right_type) {
        return false;
    }

    return runtime_value_equal(left_type, left_data, right_data);
}

bool runtime_nilinterequal(const void* left_value, const void* right_value) {
    const go_empty_interface* left;
    const go_empty_interface* right;

    left = (const go_empty_interface*)left_value;
    right = (const go_empty_interface*)right_value;
    if (left == NULL || right == NULL) {
        return left == right;
    }

    return runtime_efaceeq(left->type, left->data, right->type, right->data);
}

bool runtime_ifaceE2T2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data, void* target_value) {
    if (target_type == NULL) {
        return false;
    }

    if (target_type != source_type) {
        runtime_zero_typed_value(target_type, target_value);
        return false;
    }

    runtime_copy_typed_value(target_type, target_value, source_data);
    return true;
}

go_interface_method_table* runtime_assertitab(const go_type_descriptor* target_type, const go_type_descriptor* source_type) {
    go_interface_method_table* methods;

    if (target_type == NULL) {
        runtime_fail_simple("interface assertion has no target type");
    }
    if ((target_type->kind & GO_TYPE_KIND_MASK) != GO_TYPE_KIND_INTERFACE) {
        runtime_fail_simple("assertitab target is not an interface");
    }
    if (source_type == NULL) {
        runtime_fail_simple("interface assertion on nil value");
    }

    methods = runtime_build_interface_method_table((const go_interface_type_descriptor*)target_type, source_type);
    if (methods == NULL) {
        runtime_fail_pair("interface assertion failed", "want", runtime_pointer_value((void*)target_type), "have", runtime_pointer_value((void*)source_type));
    }

    return methods;
}

go_interface_assert_result runtime_ifaceE2I2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data) {
    go_interface_assert_result result;

    result.value.methods = NULL;
    result.value.data = NULL;
    result.ok = false;

    if (source_type == NULL) {
        return result;
    }

    result.value.methods = runtime_build_interface_method_table((const go_interface_type_descriptor*)target_type, source_type);
    if (result.value.methods == NULL) {
        return result;
    }

    result.value.data = source_data;
    result.ok = true;
    return result;
}

go_interface_assert_result runtime_ifaceI2I2(const go_type_descriptor* target_type, const go_interface_method_table* source_methods, const void* source_data) {
    const go_type_descriptor* source_type;

    source_type = NULL;
    if (source_methods != NULL) {
        source_type = source_methods->type;
    }

    return runtime_ifaceE2I2(target_type, source_type, source_data);
}

bool runtime_ifaceeq(const go_interface_method_table* left_methods, const void* left_data, const go_interface_method_table* right_methods, const void* right_data) {
    const go_type_descriptor* left_type;
    const go_type_descriptor* right_type;

    if (left_methods == NULL) {
        return right_methods == NULL;
    }
    if (right_methods == NULL) {
        return false;
    }

    left_type = left_methods->type;
    right_type = right_methods->type;
    if (left_type != right_type) {
        return false;
    }

    return runtime_value_equal(left_type, left_data, right_data);
}

bool runtime_ifacevaleq(const go_interface_method_table* left_methods, const void* left_data, const go_type_descriptor* right_type, const void* right_data) {
    const go_type_descriptor* left_type;

    if (left_methods == NULL) {
        return false;
    }

    left_type = left_methods->type;
    if (left_type != right_type) {
        return false;
    }

    return runtime_value_equal(left_type, left_data, right_data);
}

bool runtime_interequal(const void* left_value, const void* right_value) {
    const go_interface* left;
    const go_interface* right;

    if (left_value == NULL || right_value == NULL) {
        return false;
    }

    left = (const go_interface*)left_value;
    right = (const go_interface*)right_value;
    return runtime_ifaceeq(left->methods, left->data, right->methods, right->data);
}

int memcmp(const void* left, const void* right, size_t size) {
    if (left == NULL || right == NULL) {
        return left == right ? 0 : (left == NULL ? -1 : 1);
    }

    return kos_memcmp(left, right, size);
}

void* runtime_newobject(const go_type_descriptor* descriptor) {
    size_t size;
    size_t allocated_size;
    void* memory;

    size = 0;
    if (descriptor != NULL) {
        size = (size_t)descriptor->size;
    }

    allocated_size = size == 0 ? 1 : size;
    memory = malloc(allocated_size);
    if (memory == NULL) {
        return NULL;
    }

    kos_memset(memory, 0, allocated_size);
    return memory;
}

void runtime_panicmem(void) {
    runtime_fail_simple("memory or bounds failure");
}

__attribute__((noreturn)) void runtime_panicdottype(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const go_type_descriptor* interface_type) {
    (void)interface_type;

    runtime_fail_pair("type assertion failed", "want", runtime_pointer_value((void*)target_type), "have", runtime_pointer_value((void*)source_type));
}

void runtime_goPanicIndex(int32_t index, int32_t bound) {
    runtime_fail_pair("index out of range", "index", (uint32_t)index, "bound", (uint32_t)bound);
}

void runtime_goPanicIndexU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("index out of range", "index", index, "bound", bound);
}

void runtime_goPanicSliceAlen(int32_t index, int32_t bound) {
    runtime_fail_pair("slice upper bound exceeds length", "index", (uint32_t)index, "len", (uint32_t)bound);
}

void runtime_goPanicSliceAlenU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("slice upper bound exceeds length", "index", index, "len", bound);
}

void runtime_goPanicSliceAcap(int32_t index, int32_t bound) {
    runtime_fail_pair("slice upper bound exceeds capacity", "index", (uint32_t)index, "cap", (uint32_t)bound);
}

void runtime_goPanicSliceAcapU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("slice upper bound exceeds capacity", "index", index, "cap", bound);
}

void runtime_goPanicSliceB(int32_t low, int32_t high) {
    runtime_fail_pair("invalid slice bounds", "low", (uint32_t)low, "high", (uint32_t)high);
}

void runtime_goPanicSliceBU(uint32_t low, uint32_t high) {
    runtime_fail_pair("invalid slice bounds", "low", low, "high", high);
}

void runtime_goPanicSlice3Alen(int32_t index, int32_t bound) {
    runtime_fail_pair("3-index slice exceeds length", "index", (uint32_t)index, "len", (uint32_t)bound);
}

void runtime_goPanicSlice3AlenU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("3-index slice exceeds length", "index", index, "len", bound);
}

void runtime_goPanicSlice3Acap(int32_t index, int32_t bound) {
    runtime_fail_pair("3-index slice exceeds capacity", "index", (uint32_t)index, "cap", (uint32_t)bound);
}

void runtime_goPanicSlice3AcapU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("3-index slice exceeds capacity", "index", index, "cap", bound);
}

void runtime_goPanicSlice3B(int32_t index, int32_t bound) {
    runtime_fail_pair("invalid 3-index slice bounds", "index", (uint32_t)index, "bound", (uint32_t)bound);
}

void runtime_goPanicSlice3BU(uint32_t index, uint32_t bound) {
    runtime_fail_pair("invalid 3-index slice bounds", "index", index, "bound", bound);
}

void runtime_goPanicSlice3C(int32_t low, int32_t high) {
    runtime_fail_pair("invalid 3-index slice range", "low", (uint32_t)low, "high", (uint32_t)high);
}

void runtime_goPanicSlice3CU(uint32_t low, uint32_t high) {
    runtime_fail_pair("invalid 3-index slice range", "low", low, "high", high);
}

void runtime_goPanicSliceConvert(int32_t index, int32_t bound) {
    runtime_fail_pair("slice conversion out of range", "index", (uint32_t)index, "bound", (uint32_t)bound);
}

void runtime_register_gcroots(void* roots) {
    (void)roots;
}

void* memmove(void* dest, const void* src, size_t size) {
    if (dest == NULL || src == NULL) {
        return dest;
    }

    return kos_memmove(dest, src, size);
}

void* __unsafe_get_addr(void* base, size_t offset) {
    if (base == NULL) {
        return NULL;
    }

    return (void*)((unsigned char*)base + offset);
}

__asm__(".global runtime.memequal32..f");
static go_equal_function runtime_memequal32_descriptor = runtime_memequal32_impl;
__asm__(".set runtime.memequal32..f, runtime_memequal32_descriptor");

__asm__(".global runtime.memequal16..f");
static go_equal_function runtime_memequal16_descriptor = runtime_memequal16_impl;
__asm__(".set runtime.memequal16..f, runtime_memequal16_descriptor");

__asm__(".global runtime.memequal8..f");
static go_equal_function runtime_memequal8_descriptor = runtime_memequal8_impl;
__asm__(".set runtime.memequal8..f, runtime_memequal8_descriptor");

__asm__(".global runtime.memequal64..f");
static go_equal_function runtime_memequal64_descriptor = runtime_memequal64_impl;
__asm__(".set runtime.memequal64..f, runtime_memequal64_descriptor");

__asm__(".global runtime.memequal");
__asm__(".set runtime.memequal, runtime_memequal_impl");

__asm__(".global runtime.memequal64");
__asm__(".set runtime.memequal64, runtime_memequal64_impl");

__asm__(".global runtime.memequal32");
__asm__(".set runtime.memequal32, runtime_memequal32_impl");

__asm__(".global runtime.memequal16");
__asm__(".set runtime.memequal16, runtime_memequal16_impl");

__asm__(".global runtime.memequal8");
__asm__(".set runtime.memequal8, runtime_memequal8_impl");

__asm__(".global runtime.concatstrings");
__asm__(".set runtime.concatstrings, runtime_concatstrings");

__asm__(".global runtime.SetByteString");
__asm__(".set runtime.SetByteString, runtime_set_byte_string");

__asm__(".global runtime.writeBarrier");
__asm__(".set runtime.writeBarrier, runtime_write_barrier_enabled");

__asm__(".global runtime.gcWriteBarrier");
__asm__(".set runtime.gcWriteBarrier, runtime_gc_write_barrier");

__asm__(".global runtime.strequal..f");
static go_equal_function runtime_strequal_descriptor = runtime_strequal_impl;
__asm__(".set runtime.strequal..f, runtime_strequal_descriptor");

__asm__(".global runtime.strequal");
__asm__(".set runtime.strequal, runtime_strequal_impl");

__asm__(".global runtime.ifaceeq");
__asm__(".set runtime.ifaceeq, runtime_ifaceeq");

__asm__(".global runtime.ifacevaleq");
__asm__(".set runtime.ifacevaleq, runtime_ifacevaleq");

__asm__(".global runtime.efaceeq");
__asm__(".set runtime.efaceeq, runtime_efaceeq");

__asm__(".global runtime.ifaceE2T2");
__asm__(".set runtime.ifaceE2T2, runtime_ifaceE2T2");

__asm__(".global runtime.assertitab");
__asm__(".set runtime.assertitab, runtime_assertitab");

__asm__(".global runtime.ifaceE2I2");
__asm__(".set runtime.ifaceE2I2, runtime_ifaceE2I2");

__asm__(".global runtime.ifaceI2I2");
__asm__(".set runtime.ifaceI2I2, runtime_ifaceI2I2");

__asm__(".global runtime.interequal");
__asm__(".set runtime.interequal, runtime_interequal");

__asm__(".global runtime.interequal..f");
static go_equal_function runtime_interequal_descriptor = runtime_interequal;
__asm__(".set runtime.interequal..f, runtime_interequal_descriptor");

__asm__(".global runtime.nilinterequal");
__asm__(".set runtime.nilinterequal, runtime_nilinterequal");

__asm__(".global runtime.nilinterequal..f");
static go_equal_function runtime_nilinterequal_descriptor = runtime_nilinterequal;
__asm__(".set runtime.nilinterequal..f, runtime_nilinterequal_descriptor");

__asm__(".global runtime.newobject");
__asm__(".set runtime.newobject, runtime_newobject");

__asm__(".global runtime.makeslice");
__asm__(".set runtime.makeslice, runtime_makeslice");

__asm__(".global runtime.growslice");
__asm__(".set runtime.growslice, runtime_growslice");

__asm__(".global runtime.typedmemmove");
__asm__(".set runtime.typedmemmove, runtime_typedmemmove");

__asm__(".global runtime.slicebytetostring");
__asm__(".set runtime.slicebytetostring, runtime_slicebytetostring");

__asm__(".global runtime.stringtoslicebyte");
__asm__(".set runtime.stringtoslicebyte, runtime_stringtoslicebyte");

__asm__(".global runtime.memmove");
__asm__(".set runtime.memmove, memmove");

__asm__(".global runtime.panicdottype");
__asm__(".set runtime.panicdottype, runtime_panicdottype");

__asm__(".global runtime.goPanicIndex");
__asm__(".set runtime.goPanicIndex, runtime_goPanicIndex");

__asm__(".global runtime.goPanicIndexU");
__asm__(".set runtime.goPanicIndexU, runtime_goPanicIndexU");

__asm__(".global runtime.goPanicSliceAlen");
__asm__(".set runtime.goPanicSliceAlen, runtime_goPanicSliceAlen");

__asm__(".global runtime.goPanicSliceAlenU");
__asm__(".set runtime.goPanicSliceAlenU, runtime_goPanicSliceAlenU");

__asm__(".global runtime.goPanicSliceAcap");
__asm__(".set runtime.goPanicSliceAcap, runtime_goPanicSliceAcap");

__asm__(".global runtime.goPanicSliceAcapU");
__asm__(".set runtime.goPanicSliceAcapU, runtime_goPanicSliceAcapU");

__asm__(".global runtime.goPanicSliceB");
__asm__(".set runtime.goPanicSliceB, runtime_goPanicSliceB");

__asm__(".global runtime.goPanicSliceBU");
__asm__(".set runtime.goPanicSliceBU, runtime_goPanicSliceBU");

__asm__(".global runtime.goPanicSlice3Alen");
__asm__(".set runtime.goPanicSlice3Alen, runtime_goPanicSlice3Alen");

__asm__(".global runtime.goPanicSlice3AlenU");
__asm__(".set runtime.goPanicSlice3AlenU, runtime_goPanicSlice3AlenU");

__asm__(".global runtime.goPanicSlice3Acap");
__asm__(".set runtime.goPanicSlice3Acap, runtime_goPanicSlice3Acap");

__asm__(".global runtime.goPanicSlice3AcapU");
__asm__(".set runtime.goPanicSlice3AcapU, runtime_goPanicSlice3AcapU");

__asm__(".global runtime.goPanicSlice3B");
__asm__(".set runtime.goPanicSlice3B, runtime_goPanicSlice3B");

__asm__(".global runtime.goPanicSlice3BU");
__asm__(".set runtime.goPanicSlice3BU, runtime_goPanicSlice3BU");

__asm__(".global runtime.goPanicSlice3C");
__asm__(".set runtime.goPanicSlice3C, runtime_goPanicSlice3C");

__asm__(".global runtime.goPanicSlice3CU");
__asm__(".set runtime.goPanicSlice3CU, runtime_goPanicSlice3CU");

__asm__(".global runtime.goPanicSliceConvert");
__asm__(".set runtime.goPanicSliceConvert, runtime_goPanicSliceConvert");

__asm__(".global runtime.panicmem");
__asm__(".set runtime.panicmem, runtime_panicmem");

__asm__(".global runtime.registerGCRoots");
__asm__(".set runtime.registerGCRoots, runtime_register_gcroots");
