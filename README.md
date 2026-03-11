# DB Master - Taiko Wiki

태고의 달인 위키(Taiko Wiki)의 데이터베이스 관리 및 조회를 담당하는 마스터 서버입니다.

## DB 함수 레퍼런스 (`src/dbfunc/`)

모든 요청은 `/func` 엔드포인트를 통해 이루어지며, 아래는 `name`과 `params`로 전달되는 미리 정의된 함수 목록입니다.

### 1. 곡 관련 (`song.go`)

| 함수명 (Name) | 설명 | 파라미터 (Params) |
| :--- | :--- | :--- |
| `song.song-data` | 특정 곡 번호의 전체 데이터를 조회합니다. | • `songNo` (int/string) |
| `song.partial-song-data` | 특정 곡의 특정 컬럼만 선택적으로 조회합니다. | • `songNo` (int/string)<br>• `columns` ([]string) |
| `song.partial-data-of-all-songs` | 모든 곡에 대해 지정된 컬럼만 조회합니다. | • `columns` ([]string) |

### 2. 사용자 관련 (`user.go`)

| 함수명 (Name) | 설명 | 파라미터 (Params) |
| :--- | :--- | :--- |
| `user.user-data` | 제공자 정보를 통한 사용자 데이터를 조회합니다. | • `provider` (string)<br>• `providerId` (string) |
| `user.user-data-by-uuid` | UUID를 기반으로 사용자 데이터를 조회합니다. | • `uuid` (string) |
| `user.nick-and-uuid` | 모든 사용자의 닉네임과 UUID 목록을 조회합니다. | 없음 |
| `user.nicks` | 여러 UUID에 해당하는 닉네임들을 조회합니다. | • `uuids` ([]string) |
| `user.for-compe-admin` | 관리용 사용자 목록(닉네임, UUID)을 조회합니다. | 없음 |

### 3. 평점 및 프로필 관련 (`rating.go`)

| 함수명 (Name) | 설명 | 파라미터 (Params) |
| :--- | :--- | :--- |
| `rating.simple-profile` | 사용자의 요약 프로필(레이팅, 닉네임 등)을 조회합니다. | • `uuid` (string) |
| `rating.taiko-profiles` | 여러 사용자의 태고 프로필 정보를 일괄 조회합니다. | • `uuids` ([]string) |
| `rating.song-rating-datas` | 사용자의 곡별 평점 데이터를 조회합니다. (all=false 시 상위 50개) | • `uuid` (string)<br>• `all` (bool) |

### 4. 파일 로그 관련 (`file.go`)

| 함수명 (Name) | 설명 | 파라미터 (Params) |
| :--- | :--- | :--- |
| `file.getByFileName` | 파일 이름을 기반으로 업로드 로그를 조회합니다. | • `fileName` (string) |
| `file.newLog` | 새로운 파일 업로드 정보를 기록합니다. | • `UUID` (string)<br>• `originalFileName` (string)<br>• `fileName` (string) |

---

## 서버 실행 방법

1. 환경 설정 파일(`.env.json`, `connDatas.env.json`)이 준비되었는지 확인합니다.
2. `build.bat` 또는 `build.sh`를 실행하여 빌드합니다.
3. 생성된 바이너리를 실행합니다.
