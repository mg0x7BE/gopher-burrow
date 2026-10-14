//! Inventory shows the Rust elements the theme colors.

use std::collections::HashMap;
use std::fmt;
use std::sync::{Arc, Mutex};

const MAX_ITEMS: usize = 1_000;
static DEFAULT_TAGS: [&str; 2] = ["go", "theme"];

/// The state of an item.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
enum Status {
    Draft,
    Published,
    Archived,
}

#[derive(Debug, Clone)]
struct Item {
    id: u64,
    name: String,
    price: f64,
    tags: Vec<String>,
    status: Status,
}

#[derive(Debug)]
enum StoreError {
    NotFound(u64),
    Full,
}

impl fmt::Display for StoreError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            StoreError::NotFound(id) => write!(f, "item {id} not found"),
            StoreError::Full => f.write_str("store is full"),
        }
    }
}

trait Store {
    fn get(&self, id: u64) -> Result<Item, StoreError>;
    fn put(&mut self, item: Item) -> Result<(), StoreError>;
}

#[derive(Default)]
struct MemStore {
    items: HashMap<u64, Item>,
}

impl Store for MemStore {
    fn get(&self, id: u64) -> Result<Item, StoreError> {
        self.items.get(&id).cloned().ok_or(StoreError::NotFound(id))
    }

    fn put(&mut self, item: Item) -> Result<(), StoreError> {
        if self.items.len() >= MAX_ITEMS {
            return Err(StoreError::Full);
        }
        self.items.insert(item.id, item);
        Ok(())
    }
}

fn filter<'a, T>(items: &'a [T], keep: impl Fn(&T) -> bool) -> Vec<&'a T> {
    items.iter().filter(|item| keep(item)).collect()
}

fn main() {
    let mut store = MemStore::default();
    let item = Item {
        id: 42,
        name: String::from("Gopher plush"),
        price: 19.99,
        tags: DEFAULT_TAGS.iter().map(|t| t.to_string()).collect(),
        status: Status::Published,
    };

    if let Err(err) = store.put(item.clone()) {
        eprintln!("{err}");
        std::process::exit(1);
    }

    match store.get(7) {
        Ok(found) => println!("found {}", found.name),
        Err(e @ StoreError::NotFound(_)) => println!("missing item: {e}"),
        Err(e) => panic!("unexpected: {e}"),
    }

    let shared = Arc::new(Mutex::new(Vec::new()));
    let handles: Vec<_> = ["alpha", "beta"]
        .into_iter()
        .map(|name| {
            let shared = Arc::clone(&shared);
            std::thread::spawn(move || shared.lock().unwrap().push(name.to_uppercase()))
        })
        .collect();
    for handle in handles {
        handle.join().unwrap();
    }

    let raw = r"C:\path\with\backslashes";
    let (ratio, mask, letter) = (0.75, 0xFF_u8, 'g');
    let cheap = filter(std::slice::from_ref(&item), |i| i.price < 20.0);
    println!("{raw} {ratio} {mask} {letter} {} {:?}", cheap.len(), shared.lock().unwrap());
    println!("{:?} {:?} {} {}", Status::Draft, Status::Archived, item.status == Status::Published, item.tags.join(","));
}
