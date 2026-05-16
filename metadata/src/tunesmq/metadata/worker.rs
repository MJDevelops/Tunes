use crate::tunesmq::metadata::METADATA_DEALER_ADDR;
use anyhow::{Ok, Result};
use std::thread;
use zmq::{Context, Socket, SocketType};

pub struct Worker {
    thread: thread::JoinHandle<Result<()>>,
}

impl Worker {
    pub fn spawn(&mut self, context: &Context) -> Result<()> {
        let rep_sock = context.socket(SocketType::REP)?;
        rep_sock.connect(METADATA_DEALER_ADDR)?;
        self.thread = thread::spawn(move || {
            loop {}
            rep_sock.disconnect(METADATA_DEALER_ADDR)?;
        });
        Ok(())
    }
}
