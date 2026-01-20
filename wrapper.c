#include "clips.h"

void* clips_create() {
    return CreateEnvironment();
}

void clips_destroy(void* env) {
    DestroyEnvironment(env);
}

void clips_load(void* env, const char* file) {
    Load(env, file);
}

void clips_reset(void* env) {
    Reset(env);
}

void clips_run(void* env) {
    Run(env, -1);
}

void clips_assert(void* env, const char* fact) {
    AssertString(env, fact);
}
