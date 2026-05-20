use unroll::unroll;

#[test]
fn test_unroll() {
    let mut s = 0usize;
    unroll!(8, i, {
        s += i;
    });
    assert_eq!(s, 0 + 1 + 2 + 3 + 4 + 5 + 6 + 7);
}
