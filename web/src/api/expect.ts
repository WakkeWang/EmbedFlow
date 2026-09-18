export interface ExpectStep {
	await: string
	match: '' | 'contains' | 'regex' | 'exact'
	send: string
	send_before_wait: boolean
	timeout_ms: number
	on_fail: '' | 'abort' | 'continue' | 'retry'
	max_retries: number
	delay_ms: number
}

export interface ExpectRule {
	id?: number
	name: string
	steps: ExpectStep[]
}

export function emptyStep(): ExpectStep {
	return {
		await: '',
		match: '',
		send: '',
		send_before_wait: false,
		timeout_ms: 10000,
		on_fail: '',
		max_retries: 0,
		delay_ms: 0,
	}
}

// Client-side mirror of the server's expect.Validate: catch errors at the
// editor so the operator never round-trips a broken rule (issue #7).
export function validateSteps(steps: ExpectStep[]): Map<number, string> {
	const errs = new Map<number, string>()
	steps.forEach((s, i) => {
		if (!s.await && !s.send && !s.delay_ms) {
			errs.set(i, 'empty step (no await, send, or delay)')
			return
		}
		if (s.await && s.match === 'regex') {
			try {
				new RegExp(s.await)
			} catch (e) {
				errs.set(i, `bad regex: ${e}`)
			}
		}
		if (s.send && decodeSend(s.send) === null) {
			errs.set(i, 'bad send content (bad escape)')
		}
	})
	return errs
}

// decodeSend is the JS mirror of the server's EncodeSend (CEO-13A): \r \n
// \t \C \xHH. Returns null on a bad escape.
export function decodeSend(s: string): Uint8Array | null {
	const out: number[] = []
	for (let i = 0; i < s.length; i++) {
		const c = s[i]
		if (c !== '\\') {
			out.push(c.charCodeAt(0))
			continue
		}
		i++
		if (i >= s.length) return null
		switch (s[i]) {
			case 'r':
				out.push(0x0d)
				break
			case 'n':
				out.push(0x0a)
				break
			case 't':
				out.push(0x09)
				break
			case 'C':
				out.push(0x03)
				break
			case 'x': {
				if (i + 2 >= s.length) return null
				const hex = s.slice(i + 1, i + 3)
				if (!/^[0-9a-fA-F]{2}$/.test(hex)) return null
				out.push(parseInt(hex, 16))
				i += 2
				break
			}
			case '\\':
				out.push(0x5c)
				break
			default:
				return null
		}
	}
	return new Uint8Array(out)
}
