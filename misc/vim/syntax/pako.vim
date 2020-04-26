if exists("b:current_syntax")
  finish
endif

syn case match

syn keyword     pakoDirective         module
syn keyword     pakoDeclaration       var

hi def link     pakoDirective         Statement
hi def link     pakoDeclaration       Type

syn keyword     pakoStatement         return break continue throw
syn keyword     pakoConditional       if else switch try catch finally
syn keyword     pakoLabel             case default
syn keyword     pakoRepeat            for range

hi def link     pakoStatement         Statement
hi def link     pakoConditional       Conditional
hi def link     pakoLabel             Label
hi def link     pakoRepeat            Repeat

syn match       pakoDeclaration       /\<fn\>/
syn match       pakoDeclaration       /^fn\>/

syn keyword     pakoCast              bytes runes string

hi def link     pakoCast              Type

syn keyword     pakoBuiltins          keys len
syn keyword     pakoBuiltins          println printf print
syn keyword     pakoConstants         true false nil

hi def link     pakoBuiltins          Keyword
hi def link     pakoConstants         Keyword

" Comments; their contents
syn keyword     pakoTodo              contained TODO FIXME XXX BUG
syn cluster     pakoCommentGroup      contains=pakoTodo
syn region      pakoComment           start="#" end="$" contains=@pakoCommentGroup,@Spell

hi def link     pakoComment           Comment
hi def link     pakoTodo              Todo

" pako escapes
syn match       pakoEscapeOctal       display contained "\\[0-7]\{3}"
syn match       pakoEscapeC           display contained +\\[abfnrtv\\'"]+
syn match       pakoEscapeX           display contained "\\x\x\{2}"
syn match       pakoEscapeU           display contained "\\u\x\{4}"
syn match       pakoEscapeBigU        display contained "\\U\x\{8}"
syn match       pakoEscapeError       display contained +\\[^0-7xuUabfnrtv\\'"]+

hi def link     pakoEscapeOctal       pakoSpecialString
hi def link     pakoEscapeC           pakoSpecialString
hi def link     pakoEscapeX           pakoSpecialString
hi def link     pakoEscapeU           pakoSpecialString
hi def link     pakoEscapeBigU        pakoSpecialString
hi def link     pakoSpecialString     Special
hi def link     pakoEscapeError       Error

" Strings and their contents
syn cluster     pakoStringGroup       contains=pakoEscapeOctal,pakoEscapeC,pakoEscapeX,pakoEscapeU,pakoEscapeBigU,pakoEscapeError
syn region      pakoString            start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@pakoStringGroup
syn region      pakoRawString         start=+`+ end=+`+

hi def link     pakoString            String
hi def link     pakoRawString         String

" Characters; their contents
syn cluster     pakoCharacterGroup    contains=pakoEscapeOctal,pakoEscapeC,pakoEscapeX,pakoEscapeU,pakoEscapeBigU
syn region      pakoCharacter         start=+'+ skip=+\\\\\|\\'+ end=+'+ contains=@pakoCharacterGroup

hi def link     pakoCharacter         Character

" Regions
syn region      pakoBlock             start="{" end="}" transparent fold
syn region      pakoParen             start='(' end=')' transparent

" Integers
syn match       pakoDecimalInt        "\<\d\+\([Ee]\d\+\)\?\>"
syn match       pakoHexadecimalInt    "\<0x\x\+\>"
syn match       pakoOctalInt          "\<0\o\+\>"
syn match       pakoOctalError        "\<0\o*[89]\d*\>"

hi def link     pakoDecimalInt        Integer
hi def link     pakoHexadecimalInt    Integer
hi def link     pakoOctalInt          Integer
hi def link     Integer             Number

" Floating point
syn match       pakoFloat             "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syn match       pakoFloat             "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syn match       pakoFloat             "\<\d\+[Ee][-+]\d\+\>"

hi def link     pakoFloat             Float
hi def link     pakoImaginary         Number

syn sync minlines=500

let b:current_syntax = "pako"
