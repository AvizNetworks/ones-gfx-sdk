from enum import Enum


class OperationMode(str, Enum):
    SYNCHRONOUS = "sync"
    ASYNC_POLL = "async_poll"
    ASYNC_WEBHOOK = "async_webhook"


class OperationStatus(str, Enum):
    PENDING = "PENDING"
    RUNNING = "RUNNING"
    SUCCESS = "SUCCESS"
    FAILURE = "FAILURE"
    # The spelling FM actually sends for a failed operation; same meaning as FAILURE.
    FAILED = "FAILED"

    @classmethod
    def is_terminal(cls, status: str) -> bool:
        return status in (cls.SUCCESS.value, cls.FAILURE.value, cls.FAILED.value)

    @classmethod
    def is_failure(cls, status: str) -> bool:
        return status in (cls.FAILURE.value, cls.FAILED.value)


class ConfigStatus(int, Enum):
    """A device's configuration state, as FM reports it in ``config_status``.

    The numbering is positional, not a progression, and reads wrong at a
    glance::

        0  FAIL         configuration failed
        1  PASS         configuration succeeded
        2  NOT_STARTED  configuration has not begun   <- NOT a success
        3  IN_PROGRESS  configuration is underway
        4  BROWN_FIELD  skipped: the fabric deploys as brownfield
        5  SKIPPED      skipped: pre-validation or minimal-config excluded it

    The trap is ``2``. Because it sits above ``PASS`` it looks like a later,
    stronger state, so a ``status >= 1`` check silently treats a device that
    was never configured as configured. Compare against the members and never
    against the ordinal, and use :attr:`is_success` rather than ordering.

    ``4`` and ``5`` both mean "deliberately not configured" and are terminal;
    neither is a failure and neither is a success.

    FM sends these as JSON *strings* (``"1"``, not ``1``). :meth:`from_api`
    accepts either and falls back; :meth:`parse` is the strict entry point for
    values arriving from anywhere else.
    """

    FAIL = 0
    PASS = 1
    NOT_STARTED = 2
    IN_PROGRESS = 3
    BROWN_FIELD = 4
    SKIPPED = 5

    @classmethod
    def from_api(cls, value):
        """Lenient decode for a response field.

        Unknown values become ``NOT_STARTED`` rather than raising, because one
        unreadable status should not discard an otherwise good payload.
        """
        try:
            return cls(int(value))
        except (ValueError, TypeError):
            return cls.NOT_STARTED

    @classmethod
    def parse(cls, value) -> "ConfigStatus":
        """Strict decode for a value from outside a Fabric Manager response.

        Accepts a member, an int, a numeric string (``"1"``) or a canonical
        name (``"PASS"``), and raises :class:`ValueError` on anything else.
        Configuration, a CLI flag or another system gets no benefit of the
        doubt — unlike :meth:`from_api`, this fails loudly.
        """
        if isinstance(value, cls):
            return value
        if isinstance(value, bool):
            raise ValueError(f"cannot read a config status from {value!r}")
        if isinstance(value, int):
            try:
                return cls(value)
            except ValueError:
                raise ValueError(f"unknown config status {value}") from None
        if isinstance(value, str):
            s = value.strip()
            if not s:
                raise ValueError("empty config status")
            if s.lstrip("-").isdigit():
                try:
                    return cls(int(s))
                except ValueError:
                    raise ValueError(f"unknown config status {value!r}") from None
            name = s.upper().replace("-", "_")
            if name == "BROWNFIELD":
                name = "BROWN_FIELD"
            try:
                return cls[name]
            except KeyError:
                raise ValueError(f"unknown config status {value!r}") from None
        raise ValueError(f"cannot read a config status from {type(value).__name__}")

    @property
    def is_success(self) -> bool:
        """True only for ``PASS`` — the one value meaning the device is configured.

        ``BROWN_FIELD`` and ``SKIPPED`` are deliberate non-configuration, not
        success, so a caller gating on "is it configured" must not count them.
        """
        return self is ConfigStatus.PASS

    @property
    def is_terminal(self) -> bool:
        """True when the status will not change on its own.

        ``NOT_STARTED`` and ``IN_PROGRESS`` are the only non-terminal values,
        which is why polling should stop on any other — not only on ``PASS``.
        """
        return self in (
            ConfigStatus.FAIL,
            ConfigStatus.PASS,
            ConfigStatus.BROWN_FIELD,
            ConfigStatus.SKIPPED,
        )

    @property
    def is_skipped(self) -> bool:
        """True when the device was deliberately not configured."""
        return self in (ConfigStatus.BROWN_FIELD, ConfigStatus.SKIPPED)


UNLIMITED_GPUS = -1
