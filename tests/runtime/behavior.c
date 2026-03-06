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
    go_string value;
} test_box;

typedef struct {
    uint32_t words[4];
} test_raw16;

#define GO_TYPE_KIND_INTERFACE 0x14u

extern go_string runtime_concatstrings(uintptr_t ignored, const go_string* strings, size_t count);
extern void* runtime_newobject(const go_type_descriptor* descriptor);
extern void* runtime_makeslice(const go_type_descriptor* descriptor, intptr_t len, intptr_t cap);
extern go_slice runtime_growslice(const go_type_descriptor* descriptor, void* old_values, intptr_t old_len, intptr_t old_cap, intptr_t new_len);
extern void runtime_typedmemmove(const go_type_descriptor* descriptor, void* dest, const void* src);
extern go_string runtime_slicebytetostring(void* ignored, const unsigned char* src, intptr_t len);
extern go_slice runtime_stringtoslicebyte(void* ignored, const char* src, intptr_t len);
extern bool runtime_efaceeq(const go_type_descriptor* left_type, const void* left_data, const go_type_descriptor* right_type, const void* right_data);
extern bool runtime_nilinterequal(const void* left_value, const void* right_value);
extern bool runtime_ifaceE2T2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data, void* target_value);
extern go_interface_method_table* runtime_assertitab(const go_type_descriptor* target_type, const go_type_descriptor* source_type);
extern go_interface_assert_result runtime_ifaceE2I2(const go_type_descriptor* target_type, const go_type_descriptor* source_type, const void* source_data);
extern go_interface_assert_result runtime_ifaceI2I2(const go_type_descriptor* target_type, const go_interface_method_table* source_methods, const void* source_data);
extern bool runtime_interequal(const void* left_value, const void* right_value);
extern bool runtime_memequal_export(const void* left, const void* right, size_t size) __asm__("runtime.memequal");
extern void runtime_write_barrier(void** slot, void* ptr);
extern void runtime_gc_write_barrier(void** slot, void* ptr);

static int failures = 0;

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

int main(void) {
    test_concat_and_slices();
    test_allocation_and_copy();
    test_arrays_and_barriers();
    test_empty_interface_paths();
    test_interface_paths();

    if (failures != 0) {
        fprintf(stderr, "runtime behavior checks failed: %d\n", failures);
        return 1;
    }

    printf("runtime behavior checks passed\n");
    return 0;
}
