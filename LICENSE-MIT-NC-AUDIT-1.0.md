# LICENSE-MIT-NC-AUDIT-1.0.md

**Author:** projectfong  
**Copyright (c) 2025 Fong**  
**All Rights Reserved**

---

# MIT Non-Commercial Audit License (MIT-NC-AUDIT-1.0)

## 1. Permission and Scope

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to use, copy, modify, merge, publish, and distribute copies of the Software, **provided that no part of the Software or any derivative work is sold, resold, or used for direct monetary profit or commercial resale**.

For clarity:

* Commercial organizations **may use** this Software internally or for operational purposes.
* They **may not sell**, **license**, or **monetize** the Software or its derivatives.

The following conditions apply:

1. All copies or substantial portions of the Software must retain this notice and attribution.
2. The Software **shall not transmit telemetry, analytics, or usage metrics** to any external service.
3. Redistribution, modification, or reuse must include a visible link back to the original repository or author.
4. Written authorization from the Author is required for any commercial sale, resale, or monetized distribution.

---

## 2. Audit and Evidence Requirements

To preserve traceability and compliance alignment, any operational use of this Software must:

* Maintain verifiable logs of **access**, **configuration**, and **execution** events for all runtime activities.
* Store audit evidence under `/evidence/logs/` using **UTC-timestamped filenames** (for example: `requests_YYYYMMDD.log`).
* Retain integrity hashes (SHA-256 → SHA-512) for no less than:

  * **180 days** for operational logs
  * **365 days** for hash archives and verification records
* Avoid embedding any **personally identifiable information (PII)**, authentication credentials, or cryptographic secrets within logs.
* Operate under a **secure-by-default, log-everything posture**, consistent with the **ProjectFong Security Model** and its zero-trust audit principles.

---

## 3. Attribution Requirement

Any permitted fork, derivative work, or publication must include the following statement:

> “Derived from Project Fong Research Work © 2025 Fong — used under MIT-NC-AUDIT-1.0 for non-commercialization.”

A hyperlink to the original repository or author profile is required when distributed electronically.

---

## 4. Warranty Disclaimer

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, AND NON-INFRINGEMENT.
IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY CLAIM, DAMAGES, OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT, OR OTHERWISE, ARISING FROM, OUT OF, OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

## 5. Governing Law

This License shall be governed by and construed in accordance with the laws of the State of California, United States of America, without regard to its conflict of laws principles.

---

## Revision Control

| Version   | Date       | Summary                                                    | Author      |
| --------- | ---------- | ---------------------------------------------------------- | ----------- |
| **1.0.0** | 2025-10-16 | Initial MIT-NC-AUDIT publication                           | projectfong |

---


