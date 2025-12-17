fn main() {
    let mut s = String::from("hi");
    let mut r = & mut s;
    r.push_str(" means ");
    println!("r: {}", r);
    println!("s: {}", s);
}
