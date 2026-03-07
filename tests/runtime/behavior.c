#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>

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
    go_interface_method_table* methods;
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
    const go_type_descriptor common;
    const go_type_descriptor* key_type;
    const go_type_descriptor* value_type;
    const void* bucket_type;
    const void* hash_descriptor;
} go_map_type_descriptor;

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

typedef struct {
    void* value;
    bool ok;
} go_mapaccess2_result;

typedef struct {
    go_string value;
} test_box;

typedef struct {
    uint32_t words[4];
} test_raw16;

typedef struct {
    go_string label;
    intptr_t count;
} test_map_pair;

typedef struct {
    void* key;
    void* value;
    void* state;
} runtime_map_iterator;

#define GO_TYPE_KIND_INTERFACE 0x14u

extern go_string runtime_concatstrings(uintptr_t ignored, const go_string* strings, size_t count);
extern void* runtime_newobject(const go_type_descriptor* descriptor);
extern void* runtime_makeslice(const go_type_descriptor* descriptor, intptr_t len, intptr_t cap);
extern go_slice runtime_growslice(const go_type_descriptor* descriptor, void* old_values, intptr_t old_len, intptr_t old_cap, intptr_t new_len);
extern void runtime_typedmemmove(const go_type_descriptor* descriptor, void* dest, const void* src);
extern go_string runtime_slicebytetostring(void* ignored, const unsigned char* src, intptr_t len);
extern go_slice runtime_stringtoslicebyte(void* ignored, const char* src, intptr_t len);
extern bool runtime_efaceeq(const go_type_descriptor* left_type, const void* left_data, const go_type_descriptor* right_type, const void* right_data);
extern bool runtime_ifacevaleq(const go_interface_method_table* left_methods, const void* left_data, const go_type_descriptor* right_type, const void* right_data);
extern bool runtime_nilinterequal(const void* left_value, const void* right_value);
extern bool runtime_ifaceE2T2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data, void* target_value);
extern go_interface_method_table* runtime_assertitab(const go_type_descriptor* target_type, const go_type_descriptor* source_type);
extern go_interface_assert_result runtime_ifaceE2I2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data);
extern go_interface_assert_result runtime_ifaceI2I2(const go_type_descriptor* target_type, const go_interface_method_table* source_methods, const void* source_data);
extern bool runtime_interequal(const void* left_value, const void* right_value);
extern bool runtime_memequal_export(const void* left, const void* right, size_t size) __asm__("runtime.memequal");
extern void runtime_write_barrier(void** slot, void* ptr);
extern void runtime_gc_write_barrier(void** slot, void* ptr);
extern void* runtime_makemap__small(void);
extern void* runtime_makemap(const go_map_type_descriptor* map_type, intptr_t hint, void* ignored);
extern void* runtime_mapassign__fast32(const go_map_type_descriptor* map_type, void* map, uint32_t key);
extern void* runtime_mapassign__faststr(const go_map_type_descriptor* map_type, void* map, const char* key_ptr, intptr_t key_len);
extern void* runtime_mapaccess1__fast32(const go_map_type_descriptor* map_type, void* map, uint32_t key);
extern void* runtime_mapaccess1__faststr(const go_map_type_descriptor* map_type, void* map, const char* key_ptr, intptr_t key_len);
extern go_mapaccess2_result runtime_mapaccess2__fast32(const go_map_type_descriptor* map_type, void* map, uint32_t key);
extern go_mapaccess2_result runtime_mapaccess2__faststr(const go_map_type_descriptor* map_type, void* map, const char* key_ptr, intptr_t key_len);
extern void runtime_mapdelete__fast32(const go_map_type_descriptor* map_type, void* map, uint32_t key);
extern void runtime_mapdelete__faststr(const go_map_type_descriptor* map_type, void* map, const char* key_ptr, intptr_t key_len);
extern void runtime_mapiterinit(const go_map_type_descriptor* map_type, void* map, runtime_map_iterator* iterator);
extern void runtime_mapiternext(runtime_map_iterator* iterator);
extern uint32_t runtime_kos_lookup_dll_export(uint32_t table_addr, const char* name);
extern uint32_t runtime_kos_call_stdcall0(uint32_t proc);
extern uint32_t runtime_kos_call_stdcall1(uint32_t proc, uint32_t arg0);
extern uint32_t runtime_kos_call_stdcall2(uint32_t proc, uint32_t arg0, uint32_t arg1);
extern void runtime_kos_call_stdcall1_void(uint32_t proc, uint32_t arg0);
extern void runtime_kos_call_stdcall2_void(uint32_t proc, uint32_t arg0, uint32_t arg1);
extern void runtime_kos_call_stdcall5_void(uint32_t proc, uint32_t arg0, uint32_t arg1, uint32_t arg2, uint32_t arg3, uint32_t arg4);

static int failures = 0;

uint32_t runtime_kos_heap_init_raw(void) {
    return 1;
}

uint32_t runtime_kos_heap_alloc_raw(uint32_t size) {
    void* ptr = malloc(size == 0 ? 1u : (size_t)size);
    return (uint32_t)(uintptr_t)ptr;
}

uint32_t runtime_kos_heap_free_raw(uint32_t ptr) {
    free((void*)(uintptr_t)ptr);
    return 1;
}

uint32_t runtime_kos_heap_realloc_raw(uint32_t size, uint32_t ptr) {
    void* resized = realloc((void*)(uintptr_t)ptr, size == 0 ? 1u : (size_t)size);
    return (uint32_t)(uintptr_t)resized;
}

uint32_t runtime_kos_load_dll_cstring_raw(const char* path) {
    (void)path;
    return 0;
}

static bool bytes_equal(const unsigned char* left, const unsigned char* right, size_t size) {
    size_t index;

    if (left == NULL || right == NULL) {
        return left == right;
    }

    for (index = 0; index < size; index++) {
        if (left[index] != right[index]) {
            return false;
        }
    }

    return true;
}

static bool cstring_equals_len(const char* actual, const char* expected, size_t len) {
    size_t index;

    if (actual == NULL || expected == NULL) {
        return actual == expected;
    }

    for (index = 0; index < len; index++) {
        if (actual[index] != expected[index]) {
            return false;
        }
    }

    return true;
}

static void expect_true(bool value, const char* label) {
    if (!value) {
        fprintf(stderr, "FAIL: %s\n", label);
        failures++;
    }
}

static void expect_false(bool value, const char* label) {
    expect_true(!value, label);
}

static void expect_intptr_eq(intptr_t actual, intptr_t expected, const char* label) {
    if (actual != expected) {
        fprintf(stderr, "FAIL: %s (got=%ld want=%ld)\n", label, (long)actual, (long)expected);
        failures++;
    }
}

static void expect_ptr_eq(const void* actual, const void* expected, const char* label) {
    if (actual != expected) {
        fprintf(stderr, "FAIL: %s (got=%p want=%p)\n", label, actual, expected);
        failures++;
    }
}

static void expect_go_string(go_string actual, const char* expected, intptr_t expected_len, const char* label) {
    if (actual.len != expected_len) {
        fprintf(stderr, "FAIL: %s length (got=%ld want=%ld)\n", label, (long)actual.len, (long)expected_len);
        failures++;
        return;
    }
    if (!cstring_equals_len(actual.str, expected, (size_t)expected_len)) {
        fprintf(stderr, "FAIL: %s contents\n", label);
        failures++;
    }
}

static void expect_bytes(const unsigned char* actual, const unsigned char* expected, size_t size, const char* label) {
    if (!bytes_equal(actual, expected, size)) {
        fprintf(stderr, "FAIL: %s contents\n", label);
        failures++;
    }
}

static void expect_zeroed(const void* actual, size_t size, const char* label) {
    const unsigned char* bytes = (const unsigned char*)actual;
    size_t index;

    if (actual == NULL) {
        fprintf(stderr, "FAIL: %s (null)\n", label);
        failures++;
        return;
    }

    for (index = 0; index < size; index++) {
        if (bytes[index] != 0) {
            fprintf(stderr, "FAIL: %s (offset=%lu value=%u)\n", label, (unsigned long)index, (unsigned int)bytes[index]);
            failures++;
            return;
        }
    }
}

static bool test_go_string_equal(const go_string* left, const go_string* right) {
    if (left == right) {
        return true;
    }
    if (left == NULL || right == NULL) {
        return false;
    }
    if (left->len != right->len) {
        return false;
    }
    return cstring_equals_len(left->str, right->str, (size_t)left->len);
}

static bool test_string_equal(const void* left, const void* right) {
    return test_go_string_equal((const go_string*)left, (const go_string*)right);
}

static bool test_box_equal(const void* left, const void* right) {
    const test_box* left_box = (const test_box*)left;
    const test_box* right_box = (const test_box*)right;

    if (left_box == right_box) {
        return true;
    }
    if (left_box == NULL || right_box == NULL) {
        return false;
    }

    return test_go_string_equal(&left_box->value, &right_box->value);
}

static go_string test_box_label(const test_box* value) {
    return value->value;
}

static bool (*string_equal_descriptor)(const void* left, const void* right) = test_string_equal;
static bool (*box_equal_descriptor)(const void* left, const void* right) = test_box_equal;

static const go_string empty_package = {NULL, 0};
static const go_string label_name = {"Label", 5};
static const go_string other_name = {"Other", 5};
static const go_string box_name = {"box", 3};

static const go_type_descriptor label_signature_descriptor = {0};
static const go_type_descriptor string_type_descriptor = {
    sizeof(go_string),
    0,
    0,
    0,
    0,
    0,
    0,
    &string_equal_descriptor,
    NULL,
    NULL,
    NULL,
    NULL,
};
static const go_interface_method_descriptor labeler_methods[] = {
    {&label_name, &empty_package, &label_signature_descriptor},
};
static const go_interface_method_descriptor missing_methods[] = {
    {&other_name, &empty_package, &label_signature_descriptor},
};
static const go_named_type_method_descriptor box_methods[] = {
    {&label_name, &empty_package, &label_signature_descriptor, NULL, (void*)test_box_label},
};
static const go_uncommon_type box_uncommon = {
    &box_name,
    &empty_package,
    box_methods,
    1,
    1,
};
static const go_type_descriptor box_type_descriptor = {
    sizeof(test_box),
    0,
    0,
    0,
    0,
    0,
    0,
    &box_equal_descriptor,
    NULL,
    NULL,
    &box_uncommon,
    NULL,
};
static const go_interface_type_descriptor labeler_interface_descriptor = {
    {
        sizeof(go_interface),
        0,
        0,
        0,
        0,
        0,
        GO_TYPE_KIND_INTERFACE,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
    },
    labeler_methods,
    1,
    1,
};
static const go_interface_type_descriptor missing_interface_descriptor = {
    {
        sizeof(go_interface),
        0,
        0,
        0,
        0,
        0,
        GO_TYPE_KIND_INTERFACE,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
    },
    missing_methods,
    1,
    1,
};
static const go_type_descriptor raw16_descriptor = {
    sizeof(test_raw16),
    0,
    0,
    0,
    0,
    0,
    0,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
};
static const go_type_descriptor byte_descriptor = {
    1,
    0,
    0,
    0,
    0,
    0,
    0,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
};
static const go_type_descriptor int_descriptor = {
    sizeof(intptr_t),
    0,
    0,
    0,
    0,
    0,
    0,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
};
static const go_type_descriptor map_pair_descriptor = {
    sizeof(test_map_pair),
    0,
    0,
    0,
    0,
    0,
    0,
    NULL,
    NULL,
    NULL,
    NULL,
    NULL,
};
static const go_map_type_descriptor map_string_int_descriptor = {
    {
        sizeof(void*),
        0,
        0,
        0,
        0,
        0,
        0,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
    },
    &string_type_descriptor,
    &int_descriptor,
    NULL,
    NULL,
};
static const go_map_type_descriptor map_int_pair_descriptor = {
    {
        sizeof(void*),
        0,
        0,
        0,
        0,
        0,
        0,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
    },
    &int_descriptor,
    &map_pair_descriptor,
    NULL,
    NULL,
};

static void test_concat_and_slices(void) {
    const go_string parts[] = {
        {"Koli", 4},
        {NULL, 0},
        {"bri", 3},
        {"OS", 2},
    };
    const unsigned char expected_slice[] = {'K', 'O', 'S'};
    const unsigned char expected_grown[] = {'A', 'B', 'C', 0, 0, 0};
    go_string joined;
    go_slice slice;
    go_string roundtrip;
    unsigned char* made;
    go_slice grown;

    joined = runtime_concatstrings(0, parts, sizeof(parts) / sizeof(parts[0]));
    expect_go_string(joined, "KolibriOS", 9, "runtime_concatstrings joins parts");
    free((void*)joined.str);

    slice = runtime_stringtoslicebyte(NULL, "KOS", 3);
    expect_intptr_eq(slice.len, 3, "runtime_stringtoslicebyte len");
    expect_intptr_eq(slice.cap, 3, "runtime_stringtoslicebyte cap");
    expect_bytes(slice.values, expected_slice, sizeof(expected_slice), "runtime_stringtoslicebyte bytes");

    roundtrip = runtime_slicebytetostring(NULL, slice.values, slice.len);
    expect_go_string(roundtrip, "KOS", 3, "runtime_slicebytetostring roundtrip");

    made = (unsigned char*)runtime_makeslice(&byte_descriptor, 3, 5);
    expect_zeroed(made, 5, "runtime_makeslice zeroes backing storage");
    made[0] = 'A';
    made[1] = 'B';
    made[2] = 'C';

    grown = runtime_growslice(&byte_descriptor, made, 3, 3, 6);
    expect_intptr_eq(grown.len, 6, "runtime_growslice len");
    expect_intptr_eq(grown.cap, 6, "runtime_growslice cap");
    expect_bytes(grown.values, expected_grown, sizeof(expected_grown), "runtime_growslice preserves bytes");

    free(slice.values);
    free((void*)roundtrip.str);
    free(made);
    free(grown.values);
}

static void test_allocation_and_copy(void) {
    test_raw16* object;
    test_box source;
    test_box dest;

    object = (test_raw16*)runtime_newobject(&raw16_descriptor);
    expect_zeroed(object, sizeof(*object), "runtime_newobject zeroes allocated object");
    free(object);

    source.value.str = "copy";
    source.value.len = 4;
    dest.value.str = NULL;
    dest.value.len = 0;
    runtime_typedmemmove(&box_type_descriptor, &dest, &source);
    expect_go_string(dest.value, "copy", 4, "runtime_typedmemmove copies typed value");
}

static void test_arrays_and_barriers(void) {
    const unsigned char left[] = {'k', 'o', 's', '!'};
    const unsigned char equal[] = {'k', 'o', 's', '!'};
    const unsigned char different[] = {'k', 'o', 's', '?'};
    void* slot;
    uintptr_t tagged_value;

    expect_true(
        runtime_memequal_export(left, equal, sizeof(left)),
        "runtime.memequal matches equal fixed arrays");
    expect_false(
        runtime_memequal_export(left, different, sizeof(left)),
        "runtime.memequal rejects different fixed arrays");

    slot = NULL;
    tagged_value = 2026;
    runtime_write_barrier(&slot, (void*)tagged_value);
    expect_ptr_eq(slot, (void*)tagged_value, "runtime_write_barrier stores pointer");

    slot = NULL;
    runtime_gc_write_barrier(&slot, (void*)tagged_value);
    expect_ptr_eq(slot, (void*)tagged_value, "runtime_gc_write_barrier stores pointer");
}

static void test_empty_interface_paths(void) {
    go_string left;
    go_string equal_right;
    go_string different_right;
    test_box box_value;
    go_empty_interface nil_left;
    go_empty_interface nil_right;
    go_empty_interface non_nil;
    go_string asserted;
    go_string mismatch;

    left.str = "same";
    left.len = 4;
    equal_right.str = "same";
    equal_right.len = 4;
    different_right.str = "else";
    different_right.len = 4;
    box_value.value.str = "box";
    box_value.value.len = 3;

    expect_true(
        runtime_efaceeq(&string_type_descriptor, &left, &string_type_descriptor, &equal_right),
        "runtime_efaceeq matches equal strings");
    expect_false(
        runtime_efaceeq(&string_type_descriptor, &left, &string_type_descriptor, &different_right),
        "runtime_efaceeq rejects different strings");
    expect_false(
        runtime_efaceeq(&string_type_descriptor, &left, &box_type_descriptor, &box_value),
        "runtime_efaceeq rejects mismatched concrete types");

    nil_left.type = NULL;
    nil_left.data = NULL;
    nil_right.type = NULL;
    nil_right.data = NULL;
    non_nil.type = &string_type_descriptor;
    non_nil.data = &left;

    expect_true(
        runtime_nilinterequal(&nil_left, &nil_right),
        "runtime_nilinterequal matches nil empty interfaces");
    expect_false(
        runtime_nilinterequal(&nil_left, &non_nil),
        "runtime_nilinterequal rejects nil vs non-nil");

    asserted.str = NULL;
    asserted.len = 0;
    expect_true(
        runtime_ifaceE2T2(&string_type_descriptor, &string_type_descriptor, &left, &asserted),
        "runtime_ifaceE2T2 succeeds for matching concrete type");
    expect_go_string(asserted, "same", 4, "runtime_ifaceE2T2 copies asserted value");

    mismatch.str = "sentinel";
    mismatch.len = 8;
    expect_false(
        runtime_ifaceE2T2(&string_type_descriptor, &box_type_descriptor, &box_value, &mismatch),
        "runtime_ifaceE2T2 rejects mismatched concrete type");
    expect_ptr_eq(mismatch.str, NULL, "runtime_ifaceE2T2 zeroes mismatched string pointer");
    expect_intptr_eq(mismatch.len, 0, "runtime_ifaceE2T2 zeroes mismatched string length");
}

static void test_interface_paths(void) {
    test_box source;
    test_box mirror;
    test_box different;
    go_interface_method_table* methods;
    go_interface_assert_result empty_to_iface;
    go_interface_assert_result empty_to_missing;
    go_interface_method_table* source_methods;
    go_interface_assert_result iface_to_iface;
    go_interface left;
    go_interface right;
    go_interface other;
    void** table_entries;
    go_string label;

    source.value.str = "iface";
    source.value.len = 5;
    mirror.value.str = "iface";
    mirror.value.len = 5;
    different.value.str = "other";
    different.value.len = 5;

    methods = runtime_assertitab(&labeler_interface_descriptor.common, &box_type_descriptor);
    expect_true(methods != NULL, "runtime_assertitab returns IMT");
    expect_ptr_eq(methods->type, &box_type_descriptor, "runtime_assertitab stores source concrete type");
    table_entries = (void**)methods;
    expect_ptr_eq(table_entries[1], (void*)test_box_label, "runtime_assertitab stores first method pointer");
    label = ((go_string(*)(const test_box*))table_entries[1])(&source);
    expect_go_string(label, "iface", 5, "runtime_assertitab IMT dispatch works");

    empty_to_iface = runtime_ifaceE2I2(&labeler_interface_descriptor.common, &box_type_descriptor, &source);
    expect_true(empty_to_iface.ok, "runtime_ifaceE2I2 succeeds for matching interface");
    expect_true(empty_to_iface.value.methods != NULL, "runtime_ifaceE2I2 returns methods");
    expect_ptr_eq(empty_to_iface.value.data, &source, "runtime_ifaceE2I2 preserves data pointer");

    empty_to_missing = runtime_ifaceE2I2(&missing_interface_descriptor.common, &box_type_descriptor, &source);
    expect_false(empty_to_missing.ok, "runtime_ifaceE2I2 rejects missing method set");
    expect_ptr_eq(empty_to_missing.value.methods, NULL, "runtime_ifaceE2I2 mismatch leaves methods nil");
    expect_ptr_eq(empty_to_missing.value.data, NULL, "runtime_ifaceE2I2 mismatch leaves data nil");

    source_methods = runtime_assertitab(&labeler_interface_descriptor.common, &box_type_descriptor);
    iface_to_iface = runtime_ifaceI2I2(&labeler_interface_descriptor.common, source_methods, &source);
    expect_true(iface_to_iface.ok, "runtime_ifaceI2I2 succeeds for matching interface");
    expect_true(iface_to_iface.value.methods != NULL, "runtime_ifaceI2I2 returns methods");
    expect_ptr_eq(iface_to_iface.value.data, &source, "runtime_ifaceI2I2 preserves data pointer");

    left.methods = runtime_assertitab(&labeler_interface_descriptor.common, &box_type_descriptor);
    left.data = &source;
    right.methods = runtime_assertitab(&labeler_interface_descriptor.common, &box_type_descriptor);
    right.data = &mirror;
    other.methods = runtime_assertitab(&labeler_interface_descriptor.common, &box_type_descriptor);
    other.data = &different;

    expect_true(runtime_ifacevaleq(left.methods, left.data, &box_type_descriptor, &mirror), "runtime_ifacevaleq matches equal interface vs concrete value");
    expect_false(runtime_ifacevaleq(left.methods, left.data, &box_type_descriptor, &different), "runtime_ifacevaleq rejects different concrete value");
    expect_true(runtime_interequal(&left, &right), "runtime_interequal matches equal interface values");
    expect_false(runtime_interequal(&left, &other), "runtime_interequal rejects different interface values");

    free(methods);
    free(empty_to_iface.value.methods);
    free(source_methods);
    free(iface_to_iface.value.methods);
    free(left.methods);
    free(right.methods);
    free(other.methods);
}

static void test_map_paths(void) {
    void* string_map;
    intptr_t* count_slot;
    go_mapaccess2_result missing_string;
    void* int_map;
    test_map_pair* pair_slot;
    go_mapaccess2_result present_pair;
    go_mapaccess2_result missing_pair;
    runtime_map_iterator iterator;
    intptr_t sum;
    bool saw_seven;

    string_map = runtime_makemap__small();
    count_slot = (intptr_t*)runtime_mapassign__faststr(&map_string_int_descriptor, string_map, "alpha", 5);
    *count_slot = 1;
    count_slot = (intptr_t*)runtime_mapassign__faststr(&map_string_int_descriptor, string_map, "beta", 4);
    *count_slot = *(intptr_t*)runtime_mapaccess1__faststr(&map_string_int_descriptor, string_map, "alpha", 5) + 2;
    runtime_mapdelete__faststr(&map_string_int_descriptor, string_map, "alpha", 5);
    missing_string = runtime_mapaccess2__faststr(&map_string_int_descriptor, string_map, "alpha", 5);
    expect_false(missing_string.ok, "runtime map string delete clears comma-ok hit");
    expect_intptr_eq(
        *(intptr_t*)runtime_mapaccess1__faststr(&map_string_int_descriptor, string_map, "beta", 4),
        3,
        "runtime map string access returns assigned value");

    int_map = runtime_makemap(&map_int_pair_descriptor, 100, NULL);
    pair_slot = (test_map_pair*)runtime_mapassign__fast32(&map_int_pair_descriptor, int_map, 7);
    pair_slot->label.str = "seven";
    pair_slot->label.len = 5;
    pair_slot->count = 7;
    pair_slot = (test_map_pair*)runtime_mapassign__fast32(&map_int_pair_descriptor, int_map, 9);
    pair_slot->label.str = "nine";
    pair_slot->label.len = 4;
    pair_slot->count = 9;

    present_pair = runtime_mapaccess2__fast32(&map_int_pair_descriptor, int_map, 7);
    expect_true(present_pair.ok, "runtime map int access finds assigned value");
    expect_go_string(((test_map_pair*)present_pair.value)->label, "seven", 5, "runtime map int access preserves string field");
    expect_intptr_eq(((test_map_pair*)present_pair.value)->count, 7, "runtime map int access preserves int field");
    expect_intptr_eq(
        ((test_map_pair*)runtime_mapaccess1__fast32(&map_int_pair_descriptor, int_map, 9))->count,
        9,
        "runtime map int fast access returns assigned value");

    runtime_mapdelete__fast32(&map_int_pair_descriptor, int_map, 9);
    missing_pair = runtime_mapaccess2__fast32(&map_int_pair_descriptor, int_map, 9);
    expect_false(missing_pair.ok, "runtime map int delete clears comma-ok hit");

    iterator.key = NULL;
    iterator.value = NULL;
    iterator.state = NULL;
    runtime_mapiterinit(&map_int_pair_descriptor, int_map, &iterator);
    sum = 0;
    saw_seven = false;
    while (iterator.key != NULL && iterator.value != NULL) {
        uint32_t key = *(uint32_t*)iterator.key;
        test_map_pair* value = (test_map_pair*)iterator.value;
        sum += (intptr_t)key + value->count;
        if (key == 7 && value->count == 7 && cstring_equals_len(value->label.str, "seven", 5)) {
            saw_seven = true;
        }
        runtime_mapiternext(&iterator);
    }
    expect_true(saw_seven, "runtime map range exposes surviving entry");
    expect_intptr_eq(sum, 14, "runtime map range sums surviving entry");
}

#if UINTPTR_MAX == 0xFFFFFFFFu
typedef struct {
    uint32_t name;
    uint32_t data;
} test_kos_dll_export;

static uint32_t __attribute__((stdcall)) test_stdcall0_impl(void) {
    return 77;
}

static uint32_t __attribute__((stdcall)) test_stdcall1_impl(uint32_t value) {
    return value + 5;
}

static uint32_t __attribute__((stdcall)) test_stdcall2_impl(uint32_t left, uint32_t right) {
    return left ^ right;
}

static uint32_t test_void_args[5];

static void __attribute__((stdcall)) test_stdcall1_void_impl(uint32_t value) {
    test_void_args[0] = value;
}

static void __attribute__((stdcall)) test_stdcall2_void_impl(uint32_t left, uint32_t right) {
    test_void_args[0] = left;
    test_void_args[1] = right;
}

static void __attribute__((stdcall)) test_stdcall5_void_impl(uint32_t a, uint32_t b, uint32_t c, uint32_t d, uint32_t e) {
    test_void_args[0] = a;
    test_void_args[1] = b;
    test_void_args[2] = c;
    test_void_args[3] = d;
    test_void_args[4] = e;
}

static void test_kos_dll_helpers(void) {
    const char alpha[] = "alpha";
    const char beta[] = "beta";
    const char missing[] = "missing";
    test_kos_dll_export exports[] = {
        {(uint32_t)(uintptr_t)alpha, (uint32_t)(uintptr_t)test_stdcall0_impl},
        {(uint32_t)(uintptr_t)beta, (uint32_t)(uintptr_t)test_stdcall1_impl},
        {0, 0},
    };

    expect_intptr_eq(
        (intptr_t)runtime_kos_lookup_dll_export((uint32_t)(uintptr_t)exports, alpha),
        (intptr_t)(uintptr_t)test_stdcall0_impl,
        "runtime_kos_lookup_dll_export resolves first export");
    expect_intptr_eq(
        (intptr_t)runtime_kos_lookup_dll_export((uint32_t)(uintptr_t)exports, beta),
        (intptr_t)(uintptr_t)test_stdcall1_impl,
        "runtime_kos_lookup_dll_export resolves later export");
    expect_intptr_eq(
        (intptr_t)runtime_kos_lookup_dll_export((uint32_t)(uintptr_t)exports, missing),
        0,
        "runtime_kos_lookup_dll_export returns zero for missing export");

    expect_intptr_eq(
        (intptr_t)runtime_kos_call_stdcall0((uint32_t)(uintptr_t)test_stdcall0_impl),
        77,
        "runtime_kos_call_stdcall0 dispatches function pointer");
    expect_intptr_eq(
        (intptr_t)runtime_kos_call_stdcall1((uint32_t)(uintptr_t)test_stdcall1_impl, 9),
        14,
        "runtime_kos_call_stdcall1 dispatches function pointer");
    expect_intptr_eq(
        (intptr_t)runtime_kos_call_stdcall2((uint32_t)(uintptr_t)test_stdcall2_impl, 0xAA, 0x55),
        (intptr_t)(0xAA ^ 0x55),
        "runtime_kos_call_stdcall2 dispatches function pointer");

    test_void_args[0] = 0;
    runtime_kos_call_stdcall1_void((uint32_t)(uintptr_t)test_stdcall1_void_impl, 123);
    expect_intptr_eq((intptr_t)test_void_args[0], 123, "runtime_kos_call_stdcall1_void passes argument");

    test_void_args[0] = 0;
    test_void_args[1] = 0;
    runtime_kos_call_stdcall2_void((uint32_t)(uintptr_t)test_stdcall2_void_impl, 7, 11);
    expect_intptr_eq((intptr_t)test_void_args[0], 7, "runtime_kos_call_stdcall2_void passes first argument");
    expect_intptr_eq((intptr_t)test_void_args[1], 11, "runtime_kos_call_stdcall2_void passes second argument");

    test_void_args[0] = 0;
    test_void_args[1] = 0;
    test_void_args[2] = 0;
    test_void_args[3] = 0;
    test_void_args[4] = 0;
    runtime_kos_call_stdcall5_void((uint32_t)(uintptr_t)test_stdcall5_void_impl, 1, 2, 3, 4, 5);
    expect_intptr_eq((intptr_t)test_void_args[0], 1, "runtime_kos_call_stdcall5_void passes first argument");
    expect_intptr_eq((intptr_t)test_void_args[1], 2, "runtime_kos_call_stdcall5_void passes second argument");
    expect_intptr_eq((intptr_t)test_void_args[2], 3, "runtime_kos_call_stdcall5_void passes third argument");
    expect_intptr_eq((intptr_t)test_void_args[3], 4, "runtime_kos_call_stdcall5_void passes fourth argument");
    expect_intptr_eq((intptr_t)test_void_args[4], 5, "runtime_kos_call_stdcall5_void passes fifth argument");
}
#endif

int main(void) {
    test_concat_and_slices();
    test_allocation_and_copy();
    test_arrays_and_barriers();
    test_empty_interface_paths();
    test_interface_paths();
    test_map_paths();
#if UINTPTR_MAX == 0xFFFFFFFFu
    test_kos_dll_helpers();
#endif

    if (failures != 0) {
        fprintf(stderr, "runtime behavior checks failed: %d\n", failures);
        return 1;
    }

    printf("runtime behavior checks passed\n");
    return 0;
}
