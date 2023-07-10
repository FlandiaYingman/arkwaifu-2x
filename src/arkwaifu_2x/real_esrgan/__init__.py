from __future__ import annotations

from tempfile import TemporaryDirectory
from typing import TYPE_CHECKING, Final

if TYPE_CHECKING:
    from _typeshed import StrOrBytesPath

import subprocess
from enum import StrEnum
from typing import Callable

import parse
import rich.progress
from rich.console import Console
from importlib.resources import files, as_file

EXE_NAME: Final[str] = "realesrgan-ncnn-vulkan.exe"


class Model(StrEnum):
    REALESRGAN_X4PLUS = "realesrgan-x4plus"
    REALESRGAN_X4PLUS_ANIME = "realesrgan-x4plus-anime"


console = Console()


def spinner():
    return rich.progress.Progress(
        rich.progress.SpinnerColumn(),
        rich.progress.TextColumn("[progress.description]{task.description}"),
        console=console,
        transient=True,
    )


def enlarge(input_path: StrOrBytesPath,
            output_path: StrOrBytesPath,
            model: Model,
            progress: Callable[[float], None]):
    with as_file(files(__name__)) as exe_dir, TemporaryDirectory() as temp_dir:
        #   -i input-path        input image path (jpg/png/webp) or directory
        #   -o output-path       output image path (jpg/png/webp) or directory
        #   -t tile-size         tile size (>=32/0=auto, default=0) can be 0,0,0 for multi-gpu
        #   -n model-name        model name (default=realesr-animevideov3, can be realesr-animevideov3 | realesrgan-x4plus | realesrgan-x4plus-anime | realesrnet-x4plus)
        process = subprocess.Popen(
            [
                exe_dir / EXE_NAME,
                '-i', input_path,
                '-o', output_path,
                '-n', model,
            ],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            cwd=temp_dir
        )

        for line in process.stderr:
            result = parse.search("{:.2f}%", line.decode())
            percent = result.fixed[0] if result else None
            progress(percent)

        exit_code = process.wait()
        if exit_code == 0xC000013A:
            # STATUS_CONTROL_C_EXIT
            raise KeyboardInterrupt()
        if exit_code != 0:
            raise ChildProcessError(f"exit value of {process} is not 0")
