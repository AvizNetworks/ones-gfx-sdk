"""Unit tests for operation status handling.

Fabric Manager reports a failed operation as ``FAILED``; ``FAILURE`` is kept
for older payloads. Both must read as terminal failures, or a poller waits
out its whole timeout on an operation that has already ended.
"""
from __future__ import annotations

import pytest

from ones_gfx.enums import OperationStatus
from ones_gfx.models import Operation


@pytest.mark.parametrize("status", ["FAILURE", "FAILED"])
def test_failed_spellings_are_terminal_failures(status):
    assert OperationStatus.is_terminal(status)
    assert OperationStatus.is_failure(status)
    op = Operation.from_polled({"id": "op1", "status": status, "errorMessage": "boom"})
    assert op.is_done and op.is_failure and not op.is_success


@pytest.mark.parametrize("status", ["PENDING", "RUNNING"])
def test_in_flight_is_not_terminal(status):
    assert not OperationStatus.is_terminal(status)
    assert not Operation.from_polled({"id": "op1", "status": status}).is_done


def test_success_is_terminal_not_failure():
    assert OperationStatus.is_terminal("SUCCESS")
    assert not OperationStatus.is_failure("SUCCESS")
    op = Operation.from_polled({"id": "op1", "status": "SUCCESS"})
    assert op.is_done and op.is_success and not op.is_failure
