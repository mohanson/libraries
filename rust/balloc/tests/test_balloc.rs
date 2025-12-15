use balloc::{Allocator, FREE_LIST, MAX_ORDER};
use std::alloc::{GlobalAlloc, Layout};

static mut RANDSEED: u64 = 0;
pub unsafe fn rand() -> u64 {
    unsafe {
        RANDSEED = RANDSEED.wrapping_mul(6364136223846793005).wrapping_add(1);
        RANDSEED
    }
}

pub struct Knowledge {
    pub ptr: *mut u8,
    pub layout: Layout,
}

#[test]
fn test_fuzz() {
    unsafe {
        let allocator = Allocator::global();
        let mut vk: Vec<Knowledge> = vec![];
        for _ in 0..1024 * 1024 {
            let action_random = rand() as usize % 1024;
            let action = if vk.len() > action_random { 1 } else { 0 };
            match action {
                0 => {
                    let layout = Layout::from_size_align(rand() as usize % 2048, 1).unwrap();
                    let ptr = allocator.alloc(layout);
                    vk.push(Knowledge { ptr, layout });
                }
                1 => {
                    let idx = rand() as usize % vk.len();
                    let knowledge = vk.swap_remove(idx);
                    allocator.dealloc(knowledge.ptr, knowledge.layout);
                }
                _ => unreachable!(),
            }
        }
        for knowledge in vk.drain(..) {
            allocator.dealloc(knowledge.ptr, knowledge.layout);
        }
        assert_eq!(vk.len(), 0);
        for i in 0..MAX_ORDER {
            assert_eq!(FREE_LIST[i], usize::MAX);
        }
        assert_eq!(FREE_LIST[MAX_ORDER], 0);
    }
}
