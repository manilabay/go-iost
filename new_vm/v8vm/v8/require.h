#ifndef IOST_V8_REQUIRE_H
#define IOST_V8_REQUIRE_H

#include "sandbox.h"

using namespace v8;

//extern char *requireModule(SandboxPtr, const char *);
void InitRequire(Isolate *isolate, Local<ObjectTemplate> globalTpl);

#endif // IOST_V8_REQUIRE_H