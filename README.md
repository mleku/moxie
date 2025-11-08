# moxie

Systems programming with spirit.

A fork of Go with changes based on 8 years of experience of the most obstructive and bug-inducing elements of the language design, and the several limitations that make especially cryptographic operations a lot slower than C and Rust versions, and making the use of binary libraries simpler, using dynamic loading and static linking options. Eliminates the complexity and race-prone design of reference types and some built-in casting operations and explicit memory freeing.
