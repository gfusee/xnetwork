# Turns the following bash script into a python script, process should be run in parallel using asyncio.gather :

# cd /home/ubuntu/testnet/seednode
# ./seednode --log-save &
#
# cd /home/ubuntu/testnet/validator00
# ./node --operation-mode snapshotless-observer -use-log-view -log-level *:INFO --log-logger-name --log-correlation --rest-api-interface=localhost:7950 &
#
# cd /home/ubuntu/testnet/proxy
# ./proxy --log-save &

import asyncio
import os
import traceback
from pathlib import Path
from typing import Any, Coroutine, List

# ----------------- Interesting log lines -----------------


LOGLINE_GENESIS_THRESHOLD_MARKER = "started committing block"
LOGLINE_AFTER_GENESIS_INTERESTING_MARKERS = ["started committing block", "ERROR", "WARN", "DEBUG", "vm", "smartcontract"]
# We ignore SC calls on genesis.
LOGLINE_ON_GENESIS_INTERESTING_MARKERS = ["started committing block", "ERROR", "WARN", "DEBUG"]


def _is_interesting_logline(logline: str):
    is_after_genesis = False

    if LOGLINE_GENESIS_THRESHOLD_MARKER in logline:
        is_after_genesis = True

    if is_after_genesis:
        return any(e in logline for e in LOGLINE_AFTER_GENESIS_INTERESTING_MARKERS)
    return any(e in logline for e in LOGLINE_ON_GENESIS_INTERESTING_MARKERS)


def _dump_interesting_log_line(pid: int, logline: str):
    print(f"[PID={pid}]", logline)


# ----------------- Start processes -----------------

NODES_START_DELAY = 1
PROXY_START_DELAY = 10


async def run(args: List[str], cwd: Path, delay: int = 0):
    await asyncio.sleep(delay)

    process = await asyncio.create_subprocess_exec(*args, stdout=asyncio.subprocess.PIPE,
                                                   stderr=asyncio.subprocess.PIPE, cwd=cwd, limit=1024 * 512)

    pid = process.pid

    print(f"Started process [{pid}]", args)
    await asyncio.wait([
        _read_stream(process.stdout, pid),
        _read_stream(process.stderr, pid)
    ])

    return_code = await process.wait()
    print(f"Proces [{pid}] stopped. Return code: {return_code}.")


async def _read_stream(stream: Any, pid: int):
    while True:
        try:
            line = await stream.readline()
            if line:
                line = line.decode("utf-8", "replace").strip()
                if _is_interesting_logline(line):
                    _dump_interesting_log_line(pid, line)
            else:
                break
        except Exception:
            print(traceback.format_exc())


async def main():
    current_dir = os.getcwd()
    to_run: List[Coroutine[Any, Any, None]] = []

    seednode_folder = current_dir + "/seednode"
    proxy_folder = current_dir + "/proxy"

    # Seed node
    to_run.append(run(["./seednode", "--log-save"], cwd=seednode_folder))

    for i in range(4):
        validator_folder = current_dir + "/validator" + str(i).zfill(2)
        api_port = 10100 + i
        log_level = "*:DEBUG,vm:TRACE,process/smartcontract:TRACE"
        to_run.append(run([
            "./node",
            "--operation-mode=snapshotless-observer",
            "--use-log-view",
            "--log-save",
            f"--log-level={log_level}",
            "--log-logger-name",
            "--log-correlation",
            f"--rest-api-interface=localhost:{api_port}"
        ], cwd=validator_folder, delay=NODES_START_DELAY))

    # Proxy
    to_run.append(run([
        "./proxy",
        "--log-save"
    ], cwd=proxy_folder, delay=PROXY_START_DELAY))

    await asyncio.gather(*to_run)


loop = asyncio.get_event_loop()
loop.run_until_complete(main())
loop.close()
