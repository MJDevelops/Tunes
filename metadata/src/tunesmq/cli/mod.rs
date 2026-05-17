use clap::Parser;

#[derive(Parser)]
#[command(name = "tunes-meta")]
#[command(about = "Metadata service for Tunes", long_about = None)]
#[command(version = "0.1.0")]
pub(crate) struct Cli {
    /// Port number for the service to listen for messages
    #[arg(short, long)]
    pub(crate) port: u16,

    /// Number of worker threads to use
    #[arg(short, long)]
    pub(crate) threads: Option<u8>,
}
