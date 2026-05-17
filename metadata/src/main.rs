use clap::Parser;
use tunesmq::cli::Cli;

mod tunesmq;

fn main() {
    let cli = Cli::parse();

    println!("port: {:?}", cli.port);
    println!("threads: {:?}", cli.threads);
}
