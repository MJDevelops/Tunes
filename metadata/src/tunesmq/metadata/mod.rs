use std::net::TcpListener;

use crate::tunesmq::MqService;
use anyhow::{Ok, Result};
use zmq::{Context, Socket, SocketType};

mod worker;

pub(crate) const METADATA_DEALER_ADDR: &'static str = "inproc://workers";

pub struct MetadataService {
    tcp_addr: String,
    ctx: Context,
    workers: Vec<worker::Worker>,
    router_sock: Socket,
    dealer_sock: Socket,
}

impl MqService for MetadataService {
    fn start(&mut self) -> Result<()> {
        zmq::proxy(&self.router_sock, &self.dealer_sock)?;

        Ok(())
    }

    fn stop(&mut self) -> Result<()> {
        self.dealer_sock.unbind(METADATA_DEALER_ADDR)?;
        self.router_sock.unbind(&self.tcp_addr)?;

        self.ctx.destroy()?;

        Ok(())
    }
}

impl MetadataService {
    fn new(workers: u8) -> Result<Self> {
        // Find available port
        let listener = TcpListener::bind("127.0.0.1:0")?;
        let port = listener.local_addr()?.port();
        drop(listener);

        let tcp_addr = format!("tcp://127.0.0.1:{}", port);
        let ctx = Context::new();

        let router_sock = ctx.socket(SocketType::ROUTER)?;
        router_sock.bind(&tcp_addr)?;

        let dealer_sock = ctx.socket(SocketType::DEALER)?;
        dealer_sock.bind(METADATA_DEALER_ADDR)?;

        let mut threads: Vec<worker::Worker> = vec![];

        for _ in 0..workers {
            //threads.push();
        }

        Ok(Self {
            ctx,
            workers: threads,
            tcp_addr,
            router_sock,
            dealer_sock,
        })
    }
}
