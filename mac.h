#ifndef mac_h
#define mac_h

#import <Cocoa/Cocoa.h>

// This macro is used to defer the execution of a block of code in the main
// event loop.
#define defer(code)                                                            \
  dispatch_async(dispatch_get_main_queue(), ^{                                 \
                     code})

#endif /* mac_h */