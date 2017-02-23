extern crate clap;
extern crate git2;

use std::env;
use std::process::Command;
use git2::Repository;
use git2::string_array::StringArray;
use clap::{Arg, App};

fn print_remotes(remotes: StringArray) {
    for remote in remotes.iter() {
        match remote {
            Some(rem) => println!("{}", rem),
            None => println!("Shitbox"),
        }
    }
}

fn tests() {
    let target_commit = env::args().nth(1).expect("Missing target commit.");

    let repo = match Repository::open(".") {
        Ok(repo) => repo,
        Err(e) => panic!("failed to open: {}", e),
    };

    println!("{:?}", repo.state());

    match repo.remotes() {
        Ok(rem) => print_remotes(rem),
        Err(e) => panic!("Oh shit! {}", e),
    };

    match repo.revparse_single(&target_commit) {
        Ok(spec) => println!("{:?}", &spec.id()),
        Err(e) => panic!("Invalid revspec {}: {}!", target_commit, e),
    };

    let output = Command::new("git")
        .arg("rev-parse")
        .arg(&target_commit)
        .output()
        .expect("failed to execute process");

    println!("status: {}", &output.status);
    println!("stdout: {}", String::from_utf8_lossy(&output.stdout));
    println!("stderr: {}", String::from_utf8_lossy(&output.stderr));
}

fn main() {
    let matches = App::new("git-absorb")
        .version("0.1.0")
        .author("Axel von Bertoldi <bertoldia@gmail.com>")
        .about("Absorb (i.e. merge) outstanding changes into the specified commit.")
        .arg(Arg::with_name("print")
            .short("p")
            .long("print-candidates")
            .conflicts_with("target")
            .help("Print the SHA1 of all candidate commits i.e. all commits in the branch in a \
                   human readable format."))
        .arg(Arg::with_name("machine")
            .short("m")
            .long("machine-parsable")
            .requires("print")
            .help("Print the candidate commits in a machine-parsable format (i.e. just the full \
                   SHA1)."))
        .arg(Arg::with_name("target")
            .index(1)
            .conflicts_with("print")
            .required(true)
            .help("Print the candidate commits in a machine-parsable format (i.e. just the full \
                   SHA1)."))
        .get_matches();
}
