# THIRD_PARTY_LICENSES.md

**Project:** cfo-scratchpad
**Author:** projectfong
**License:** MIT-NC-AUDIT-1.0
**Last Updated:** 2025-10-16

---

## 1. Purpose Summary

To document any third-party components, libraries, or tools that may be included, referenced, or integrated within this project.
This file supports audit traceability and compliance alignment under the ProjectFong Security Model and MIT-NC-AUDIT-1.0 License.

---

## 2. Policy Statement

All third-party code, binaries, or external dependencies used by this project must:

* Be **openly licensed** (permissive or equivalent).
* Contain **no embedded telemetry, tracking, or analytics**.
* Comply with **non-commercial redistribution** under MIT-NC-AUDIT-1.0.
* Be clearly listed with **license attribution** and **source reference** below.

No closed-source, proprietary, or telemetry-enabled components may be linked or embedded.

---

## 3. Verified Third-Party Components

| Component                   | Purpose                                         | License        | Source / URL                                                                                   |
| --------------------------- | ----------------------------------------------- | -------------- | ---------------------------------------------------------------------------------------------- |
| **Go Standard Library**     | Backend HTTP, file handling, and JSON encoding  | BSD-style (Go) | [https://golang.org/LICENSE](https://golang.org/LICENSE)                                       |
| **BusyBox (optional)**      | Lightweight shell utilities for Docker runtime  | GPLv2          | [https://busybox.net/license.html](https://busybox.net/license.html)                           |
| **Docker Engine / Compose** | Container orchestration for reproducible builds | Apache 2.0     | [https://www.docker.com/legal/open-source](https://www.docker.com/legal/open-source)           |
| **Mermaid (Docs only)**     | Diagram rendering for documentation             | MIT            | [https://github.com/mermaid-js/mermaid](https://github.com/mermaid-js/mermaid)                 |
| **Node.js (optional)**      | Local static asset tooling (if used)            | MIT            | [https://nodejs.org/en/about/resources/license](https://nodejs.org/en/about/resources/license) |
| **Linux Core Utilities**    | Base system and cron for evidence rotation      | GPLv3          | [https://www.gnu.org/licenses/gpl-3.0.html](https://www.gnu.org/licenses/gpl-3.0.html)         |

> **Note:** The `frontend/` directory uses plain HTML, JS, and CSS.
> No React, Vue, Angular, or telemetry-emitting frameworks are included.

---

## 4. Non-Included Components

* No proprietary SDKs, analytics scripts, or cloud telemetry agents.
* No embedded third-party API keys or tracking endpoints.
* No AI models, external network calls, or auto-updaters.

---

## 5. Compliance Assurance

All included components were reviewed for:

* License compatibility with MIT-NC-AUDIT-1.0
* Absence of telemetry and third-party data sharing
* Minimal and auditable footprint suitable for air-gapped environments

---

## 6. Revision Control

| Version   | Date       | Summary                                                  | Author      |
| --------- | ---------- | -------------------------------------------------------- | ----------- |
| **1.0.0** | 2025-10-16 | Initial publication of THIRD-PARTY.md for cfo-scratchpad | projectfong |

