import sys
from pathlib import Path

import pytest

# Run against the source tree without requiring `pip install -e .`.
sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from ones_gfx import Client  # noqa: E402

from .fake_fm import TOKEN, FakeFabricManager  # noqa: E402


@pytest.fixture
def fm():
    fake = FakeFabricManager()
    fake.base_url = fake.start()
    yield fake
    fake.stop()


@pytest.fixture
def live_client(fm):
    """A real Client talking HTTP to the fake Fabric Manager."""
    with Client.initialize_with_token(fm.base_url, TOKEN) as client:
        yield client
