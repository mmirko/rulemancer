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

char *find_facts_as_string(void *env, const char *fact_name) {
    Fact *fact = GetNextFact(env, NULL);

    StringBuilder *sb = CreateStringBuilder(env, 1024);
    if (!sb) return NULL;

    while (fact != NULL) {
        CLIPSLexeme *rel = FactRelation(fact);
        const char *relation = rel ? rel->contents : NULL;

        if (relation && strcmp(relation, fact_name) == 0) {
            FactPPForm(fact, sb, false);
            AppendStrings(env, sb->contents,"");
        }

        fact = GetNextFact(env,fact);
    }

    char *result = CopyString(env, sb->contents);

    SBDispose(sb);
    return result;
}

char * find_all_facts_as_string(void *env) {
    Fact *fact = GetNextFact(env, NULL);

    StringBuilder *sb = CreateStringBuilder(env, 1024);
    if (!sb) return NULL;

    while (fact != NULL) {
        FactPPForm(fact, sb, false);
        AppendStrings(env, sb->contents,"\n");

        fact = GetNextFact(env,fact);
    }

    char *result = CopyString(env, sb->contents);

    SBDispose(sb);
    return result;
}

void clips_free_string(void *env,char *str) {
    if (str != NULL) {
        rm(env, str, strlen(str)+1);
    }
}