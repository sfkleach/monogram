compile_mode :pop11 +strict;

section;

;;; Idiom for 'silently load a lib defining a class if not loaded already'.
if identprops("options_key") then
    loadlib("options")
endif;



;;;
;;; Pop-11 really does not like the identifier dollar_{ so we need to force
;;; the assignment with some low-level code.
;;;
ident_declare( "'dollar_{'", "syntax", 0 );
procedure(word) with_props 'dollar_{';
    $-options$-compile_newoptions_to( "}" ) -> _;
    ;;;compile_newoptions_to( "}" ) -> _;
endprocedure -> idval( identof( "'dollar_{'" ) );

endsection;
