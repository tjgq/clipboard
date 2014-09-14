#import <Cocoa/Cocoa.h>
#include <string.h>

long count() {
  NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
  return [pasteboard changeCount];
}

char *get() {
  NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
  NSArray *classes = [NSArray arrayWithObject:[NSString class]];
  NSDictionary *options = [NSDictionary dictionary];

  NSArray *items = [pasteboard readObjectsForClasses:classes options:options];
  if (items != nil && [items count] > 0) {
    return strdup([[items objectAtIndex:0] UTF8String]);
  } else {
    return NULL;
  }
}

int set(const char *s) {
  NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
  NSString *string = [NSString stringWithUTF8String:s];
  NSArray *items = [NSArray arrayWithObject:string];
  [pasteboard clearContents];
  return [pasteboard writeObjects:items];
}
