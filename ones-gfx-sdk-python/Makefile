# ----------------------------------------------------------------------------
# ONES Spectrum-X SDK — developer convenience targets.
#
# Usage:
#     make help           # list all targets with descriptions
#     make install        # editable install into current Python env
#     make build          # build wheel + sdist into dist/
#     make run-example    # run examples/usage_example.py
#     make clean          # remove build artifacts and caches
#
# All targets assume the active ``python`` and ``pip`` belong to the
# environment you want to install into. Activate a venv first if needed.
# ----------------------------------------------------------------------------

# Use the active python by default; override with `make PYTHON=python3.11 ...`
PYTHON ?= python
PIP    ?= $(PYTHON) -m pip

# Where the wheel lands.
DIST_DIR := dist

# Default target shows help if you run plain `make`.
.DEFAULT_GOAL := help

# All phony — none of these targets produce a file matching the target name.
.PHONY: help install install-dev build wheel sdist run-example \
        check clean uninstall reinstall venv

# ----------------------------------------------------------------------------
# help — auto-generated from the `## ` comments after each target.
# ----------------------------------------------------------------------------
help:  ## Show this help message
	@echo "ONES Spectrum-X SDK — make targets"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Override Python with e.g.:  make PYTHON=python3.11 install"

# ----------------------------------------------------------------------------
# Install / uninstall.
# ----------------------------------------------------------------------------
install:  ## Editable install of the SDK (pip install -e .)
	$(PIP) install -e .

install-dev:  ## Editable install plus dev extras (pytest, ruff)
	$(PIP) install -e ".[dev]"

uninstall:  ## Remove the SDK from the current Python environment
	$(PIP) uninstall -y ones-gfx-sdk

reinstall: uninstall install  ## Uninstall then install (useful after big changes)

# ----------------------------------------------------------------------------
# Build distributable artifacts.
# ----------------------------------------------------------------------------
build:  ## Build wheel + sdist into dist/
	$(PIP) install --quiet build
	$(PYTHON) -m build

wheel:  ## Build only the wheel into dist/
	$(PIP) install --quiet build
	$(PYTHON) -m build --wheel

sdist:  ## Build only the source distribution into dist/
	$(PIP) install --quiet build
	$(PYTHON) -m build --sdist

# ----------------------------------------------------------------------------
# Run the bundled usage example. No install required — the example
# bootstraps sys.path itself.
# ----------------------------------------------------------------------------
run-example:  ## Run examples/usage_example.py against your configured ONES instance
	$(PYTHON) examples/usage_example.py

# ----------------------------------------------------------------------------
# Quick smoke check that the package imports without errors.
# Useful in CI as a minimal sanity gate before running real tests.
# ----------------------------------------------------------------------------
check:  ## Quick import check — verifies the package loads cleanly
	$(PYTHON) -c "import ones_gfx; print('ones_gfx', ones_gfx.__version__, 'OK')"

# ----------------------------------------------------------------------------
# Cleanup — remove build outputs, caches, and bytecode.
# ----------------------------------------------------------------------------
clean:  ## Remove build artifacts, caches, and *.pyc files
	rm -rf $(DIST_DIR) build *.egg-info
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type d -name .pytest_cache -exec rm -rf {} +
	find . -type d -name .ruff_cache -exec rm -rf {} +
	find . -type f -name '*.pyc' -delete

# ----------------------------------------------------------------------------
# Optional: spin up a local virtualenv at .venv. After running this,
# activate it manually:  source .venv/bin/activate
# ----------------------------------------------------------------------------
venv:  ## Create a local .venv (you still need to activate it manually)
	$(PYTHON) -m venv .venv
	@echo ""
	@echo "Activate with:  source .venv/bin/activate    (Linux/macOS)"
	@echo "                .venv\\Scripts\\activate       (Windows)"
