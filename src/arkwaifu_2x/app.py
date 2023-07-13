from __future__ import annotations

import os
from concurrent.futures import as_completed
from concurrent.futures.thread import ThreadPoolExecutor
from tempfile import NamedTemporaryFile
from typing import TYPE_CHECKING, TypeAlias, Callable

import inquirer
from PIL import Image
from rich.console import Console, Group
from rich.live import Live
from rich.progress import Progress, SpinnerColumn, BarColumn, MofNCompleteColumn, TaskProgressColumn, TimeElapsedColumn, \
    TimeRemainingColumn

from arkwaifu_2x import arkwaifu, real_esrgan
from arkwaifu_2x.real_esrgan import Model

if TYPE_CHECKING:
    from _typeshed import StrOrBytesPath

console = Console()

ProgressUpdater: TypeAlias = Callable[[float], None]


def create_total_progress():
    return Progress(
        SpinnerColumn(),
        "[progress.description]{task.description}",
        BarColumn(),
        MofNCompleteColumn(),
        TaskProgressColumn(),
        TimeElapsedColumn(),
        TimeRemainingColumn(),
    )


def create_sub_progress():
    return Progress(
        "  [progress.description]{task.description}",
        BarColumn(),
        TaskProgressColumn(),
        TimeElapsedColumn(),
        TimeRemainingColumn(),
    )


def fetch_art(id: str) -> StrOrBytesPath:
    b = arkwaifu.get_content(id, "origin")
    with NamedTemporaryFile(suffix=".webp", delete=False) as temporary_file:
        temporary_file.write(b)
        return temporary_file.name


def enlarge_art(path: StrOrBytesPath, model: Model, updater: ProgressUpdater) -> StrOrBytesPath:
    try:
        temporary_file = NamedTemporaryFile(suffix=".png", delete=False)
        temporary_file.close()

        real_esrgan.enlarge(
            path,
            temporary_file.name,
            model,
            updater
        )
        return temporary_file.name
    finally:
        os.remove(path)


def convert_art(path: StrOrBytesPath) -> StrOrBytesPath:
    try:
        with Image.open(path) as image, NamedTemporaryFile(suffix=".webp", delete=False) as temporary_file:
            image.save(temporary_file, quality=90)
            return temporary_file.name
    finally:
        os.remove(path)


def submit_art(path: StrOrBytesPath, id: str, variation: str):
    try:
        arkwaifu.put_variant(id, variation)
        arkwaifu.put_content(id, variation, path)
    finally:
        os.remove(path)


def enlarge_arts(variation: str, model: real_esrgan.Model):
    with console.status(f"Finding the arts that have to be enlarged to [cyan]{variation}[/]... "):
        arts = arkwaifu.get_arts_by_absent_variation(variation)

    if len(arts) > 0:
        console.log(f"Found {len(arts)} arts that have to be enlarged to [cyan]{variation}[/].")
    else:
        console.log(f"No arts that have to be enlarged to [cyan]{variation}[/]. Exit.")
        return

    total_progress = create_total_progress()
    sub_progress = create_sub_progress()

    console.log(f"Start to enlarging {len(arts)} arts to [cyan]{variation}[/]. ")
    console.log(f"Note: the ETA may be inaccurate. ")
    with Live(Group(total_progress, sub_progress), console=console, transient=True):
        executor = ThreadPoolExecutor(max_workers=2)

        def worker(id: str):
            sub_task_id = sub_progress.add_task("")
            try:
                def progress(percentage: float):
                    sub_progress.update(sub_task_id, completed=percentage)

                sub_progress.update(sub_task_id, description=f"Fetching art: [cyan]{id}[/]... ")
                progress(0.0)
                fetched_art = fetch_art(id)

                sub_progress.update(sub_task_id, description=f"Enlarging art: [cyan]{id}[/]... ", )
                progress(0.0)
                enlarged_art = enlarge_art(fetched_art, model, progress)

                sub_progress.update(sub_task_id, description=f"Converting art: [cyan]{id}[/]... ", )
                progress(0.0)
                converted_art = convert_art(enlarged_art)

                sub_progress.update(sub_task_id, description=f"Submitting art: [cyan]{id}[/]... ", )
                progress(0.0)
                submit_art(converted_art, id, variation)
            except KeyboardInterrupt:
                executor.shutdown(wait=False, cancel_futures=True)
            finally:
                sub_progress.remove_task(sub_task_id)

        try:
            total_task_id = total_progress.add_task(f"Enlarging arts to [cyan]{variation}[/]... ", total=len(arts))
            futures = [executor.submit(worker, art.id) for art in arts]
            for future in as_completed(futures):
                try:
                    future.result()
                except Exception:
                    console.print_exception(width=None)
                finally:
                    total_progress.advance(total_task_id)
        except KeyboardInterrupt:
            console.log("[red]Keyboard Interrupt: Waiting for current progresses to complete... [/]")
        finally:
            executor.shutdown(cancel_futures=True, wait=True)


def main():
    choices = [
        "1. Enlarge the Arts with real-esrgan(realesrgan-x4plus).",
        "2. Enlarge the Arts with real-esrgan(realesrgan-x4plus-anime)."
    ]
    choice = inquirer.list_input("What would you like to do with Arkwaifu 2x? ", choices=choices)
    match choices.index(choice):
        case 0:
            enlarge_arts("real-esrgan(realesrgan-x4plus)", Model.REALESRGAN_X4PLUS)
        case 1:
            enlarge_arts("real-esrgan(realesrgan-x4plus-anime)", Model.REALESRGAN_X4PLUS_ANIME)
